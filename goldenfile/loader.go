package goldenfile

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/internal/rootfs"
	"github.com/dogmatiq/aureus/test"
)

// Loader loads [test.Test] values from directories containing pairs of files
// representing input and expected output.
type Loader struct {
	options loadOptions
}

// NewLoader returns a new [Loader], which loads golden file tests from the
// filesystem.
func NewLoader(options ...LoadOption) *Loader {
	l := &Loader{
		options: loadOptions{
			FileSystem:   rootfs.FS,
			Recursive:    true,
			IsGoldenFile: DefaultPredicate,
		},
	}
	for _, opt := range options {
		opt(&l.options)
	}
	return l
}

// Load returns tests based on the golden files in the given directory.
func (l *Loader) Load(dir string, options ...LoadOption) (test.Test, error) {
	opts := l.options
	for _, opt := range options {
		opt(&opts)
	}
	return loadDir(opts, dir, test.EmptyFlagSet)
}

func loadDir(
	opts loadOptions,
	dirPath string,
	inherited test.FlagSet,
) (test.Test, error) {
	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	t := test.Test{
		Name:   name,
		Flags:  inherited,
		Origin: test.DirectoryOrigin{DirPath: dirPath},
	}

	if skip {
		t.Flags.Add(test.FlagSkipped)
		inherited.Add(test.FlagAncestorSkipped)
	}

	entries, err := fs.ReadDir(opts.FileSystem, dirPath)
	if err != nil {
		return test.Test{}, err
	}

	var loaders []func() (test.Test, error)
	var inputs []string

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}

		n := path.Join(dirPath, e.Name())

		if e.IsDir() {
			if opts.Recursive {
				sub, err := loadDir(opts, n, inherited)
				if err != nil {
					return test.Test{}, err
				}
				t.SubTests = append(t.SubTests, sub)
			}
		} else if isInputFile, ok := opts.IsGoldenFile(e.Name()); ok {
			loaders = append(
				loaders,
				func() (test.Test, error) {
					var ins []string
					for _, filename := range inputs {
						if isInputFile(path.Base(filename)) {
							ins = append(ins, filename)
						}
					}
					return loadGoldenFile(opts, n, ins, inherited)
				},
			)
		} else {
			inputs = append(inputs, n)
		}
	}

	for _, load := range loaders {
		sub, err := load()
		if err != nil {
			return test.Test{}, err
		}
		t.SubTests = append(t.SubTests, sub)
	}

	return t, nil
}

func loadGoldenFile(
	opts loadOptions,
	filePath string,
	inputs []string,
	inherited test.FlagSet,
) (test.Test, error) {
	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")

	t := test.Test{
		Name:   name,
		Flags:  inherited,
		Origin: test.FileOrigin{FilePath: filePath},
	}

	if skip {
		t.Flags.Add(test.FlagSkipped)
		inherited.Add(test.FlagAncestorSkipped)
	}

	output, err := loadContent(opts, filePath)
	if err != nil {
		return test.Test{}, err
	}

	for _, filePath := range inputs {
		sub, err := loadInputFile(opts, filePath, inherited, output)
		if err != nil {
			return test.Test{}, err
		}
		t.SubTests = append(t.SubTests, sub)
	}

	return t, nil
}

func loadInputFile(
	opts loadOptions,
	filePath string,
	inherited test.FlagSet,
	output test.Content,
) (test.Test, error) {
	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")

	input, err := loadContent(opts, filePath)
	if err != nil {
		return test.Test{}, err
	}

	t := test.Test{
		Name:   name,
		Flags:  inherited,
		Origin: input.Origin,
		Assertions: []test.Assertion{
			test.EqualAssertion{
				Input:  input,
				Output: output,
			},
		},
	}

	if skip {
		t.Flags.Add(test.FlagSkipped)
	}

	return t, nil
}

// loadContent loads the content of a golden file or input file.
func loadContent(
	opts loadOptions,
	filePath string,
) (test.Content, error) {
	data, err := fs.ReadFile(opts.FileSystem, filePath)
	if err != nil {
		return test.Content{}, err
	}

	return test.Content{
		Origin: test.FileOrigin{FilePath: filePath},
		Data:   string(data),
	}, nil
}
