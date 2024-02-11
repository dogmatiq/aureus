package goldenfile

import (
	"fmt"
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
			FS:       rootfs.FS,
			Recurse:  true,
			IsOutput: IsOutputFile,
		},
	}

	for _, opt := range options {
		opt(&l.options)
	}

	return l
}

// Load returns a test built from files in the given directory.
//
// Any directory or file that begins with an underscore produces a test that is
// marked as skipped.
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

	parent, inherited := test.New(
		test.WithName(name),
		test.If(skip, test.WithFlag(test.FlagSkipped)),
		test.WithInheritedFlags(inherited),
	)

	entries, err := fs.ReadDir(opts.FS, dirPath)
	if err != nil {
		return test.Test{}, err
	}

	var (
		inputs  []*inputFile
		outputs []*outputFile
	)

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}

		n := path.Join(dirPath, e.Name())

		if e.IsDir() {
			if opts.Recurse {
				child, err := loadDir(opts, n, inherited)
				if err != nil {
					return test.Test{}, err
				}
				parent.SubTests = append(parent.SubTests, child)
			}
		} else if isInput, ok := opts.IsOutput(e.Name()); ok {
			out := &outputFile{Path: n, IsInput: isInput}
			outputs = append(outputs, out)
			for _, in := range inputs {
				if err := correlate(opts, in, out); err != nil {
					return test.Test{}, err
				}
			}
		} else {
			in := &inputFile{Path: n}
			inputs = append(inputs, in)
			for _, out := range outputs {
				if err := correlate(opts, in, out); err != nil {
					return test.Test{}, err
				}
			}
		}
	}

	for _, out := range outputs {
		if len(out.Inputs) == 0 {
			return test.Test{}, fmt.Errorf("output file %q has no associated input files", out.Path)
		}

		for _, in := range out.Inputs {
			child, err := buildTest(opts, in, out, inherited)
			if err != nil {
				return test.Test{}, err
			}

			parent.SubTests = append(parent.SubTests, child)
		}
	}

	for _, in := range inputs {
		if len(in.Outputs) == 0 {
			return test.Test{}, fmt.Errorf("input file %q has no associated output files", in.Path)
		}
	}

	return parent, nil
}

func buildTest(
	opts loadOptions,
	in *inputFile,
	out *outputFile,
	inherited test.FlagSet,
) (test.Test, error) {
	outName, skipOut := parsePath(out.Path)
	inName, skipIn := parsePath(in.Path)

	name := outName
	if outName != inName {
		name = fmt.Sprintf("%s -> %s", inName, outName)
	}

	t, _ := test.New(
		test.WithName(name),
		test.If(skipIn || skipOut, test.WithFlag(test.FlagSkipped)),
		test.WithInheritedFlags(inherited),
		test.WithAssertion(
			test.EqualAssertion{
				Input:  *in.Content,
				Output: *out.Content,
			},
		),
	)

	return t, nil
}

// loadContent loads the content of an input file or output file.
func loadContent(
	opts loadOptions,
	filePath string,
) (*test.Content, error) {
	data, err := fs.ReadFile(opts.FS, filePath)
	if err != nil {
		return nil, err
	}

	return &test.Content{
		Data: string(data),
		File: filePath,
	}, nil
}

type outputFile struct {
	Path           string
	IsInput        InputPredicate
	IsMatrixMember bool
	Content        *test.Content
	Inputs         []*inputFile
}

type inputFile struct {
	Path           string
	IsMatrixMember bool
	Content        *test.Content
	Outputs        []*outputFile
}

func correlate(opts loadOptions, in *inputFile, out *outputFile) error {
	if !out.IsInput(path.Base(in.Path)) {
		return nil
	}

	if out.Content == nil {
		var err error
		out.Content, err = loadContent(opts, out.Path)
		if err != nil {
			return err
		}
	}

	if in.Content == nil {
		var err error
		in.Content, err = loadContent(opts, in.Path)
		if err != nil {
			return err
		}
	}

	in.Outputs = append(in.Outputs, out)
	out.Inputs = append(out.Inputs, in)

	if len(in.Outputs) > 1 || len(out.Inputs) > 1 {
		in.IsMatrixMember = true
		out.IsMatrixMember = true
	}

	return nil
}

func parsePath(p string) (string, bool) {
	name := path.Base(p)
	name = strings.Split(name, ".")[0]
	return strings.CutPrefix(name, "_")
}
