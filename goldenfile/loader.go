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
	prototype loader
}

// LoadOption is an option that changes the behavior of a [Loader].
type LoadOption func(*loader)

// WithRecursion if a [LoadOption] that enables or disables recursive scanning
// of sub-directories.
//
// Recursion is enabled by default.
func WithRecursion(on bool) LoadOption {
	return func(l *loader) {
		l.isRecursive = on
	}
}

// NewLoader returns a new [Loader], which loads golden file tests from the
// filesystem.
func NewLoader(options ...LoadOption) *Loader {
	l := &Loader{
		prototype: loader{
			filesytem:    rootfs.FS,
			isRecursive:  true,
			isGoldenFile: DefaultPredicate,
		},
	}
	for _, opt := range options {
		opt(&l.prototype)
	}
	return l
}

// Load returns tests based on the golden files in the given directory.
func (l *Loader) Load(dir string, options ...LoadOption) (test.Runnable, error) {
	loader := l.prototype
	for _, opt := range options {
		opt(&loader)
	}
	return loader.loadDir(dir, test.EmptyFlagSet)
}

type loader struct {
	filesytem    fs.FS
	isRecursive  bool
	isGoldenFile Predicate
}

func (l *loader) loadDir(dirPath string, inherited test.FlagSet) (test.Runnable, error) {
	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	suite := &test.Suite{
		Name:   name,
		Flags:  inherited,
		Origin: test.DirectoryOrigin{DirPath: dirPath},
	}

	if skip {
		suite.Flags.Add(test.FlagSkipped)
		inherited.Add(test.FlagAncestorSkipped)
	}

	entries, err := fs.ReadDir(l.filesytem, dirPath)
	if err != nil {
		return nil, err
	}

	var loaders []func() (test.Runnable, error)
	var inputs []string

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}

		n := path.Join(dirPath, e.Name())

		if e.IsDir() {
			if l.isRecursive {
				sub, err := l.loadDir(n, inherited)
				if err != nil {
					return nil, err
				}
				suite.Tests = append(suite.Tests, sub)
			}
		} else if isInputFile, ok := l.isGoldenFile(e.Name()); ok {
			loaders = append(
				loaders,
				func() (test.Runnable, error) {
					var ins []string
					for _, filename := range inputs {
						if isInputFile(path.Base(filename)) {
							ins = append(ins, filename)
						}
					}
					return l.loadGoldenFile(n, ins, inherited)
				},
			)
		} else {
			inputs = append(inputs, n)
		}
	}

	for _, load := range loaders {
		sub, err := load()
		if err != nil {
			return nil, err
		}
		suite.Tests = append(suite.Tests, sub)
	}

	return suite, nil
}

func (l *loader) loadGoldenFile(
	filePath string,
	inputs []string,
	inherited test.FlagSet,
) (test.Runnable, error) {
	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")

	suite := &test.Suite{
		Name:   name,
		Flags:  inherited,
		Origin: test.FileOrigin{FilePath: filePath},
	}

	if skip {
		suite.Flags.Add(test.FlagSkipped)
		inherited.Add(test.FlagAncestorSkipped)
	}

	ouptut, err := l.loadContent(filePath)
	if err != nil {
		return nil, err
	}

	for _, filePath := range inputs {
		sub, err := l.loadInputFile(filePath, inherited, ouptut)
		if err != nil {
			return nil, err
		}
		suite.Tests = append(suite.Tests, sub)
	}

	return suite, nil
}

func (l *loader) loadInputFile(
	filePath string,
	inherited test.FlagSet,
	output test.Content,
) (test.Runnable, error) {
	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")

	input, err := l.loadContent(filePath)
	if err != nil {
		return nil, err
	}

	t := &test.Test{
		Name:   name,
		Flags:  inherited,
		Origin: input.Origin,
		Assertions: []test.Assertion{
			&test.Equal{
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

func (l *loader) loadContent(filePath string) (test.Content, error) {
	data, err := fs.ReadFile(l.filesytem, filePath)
	if err != nil {
		return test.Content{}, err
	}

	return test.Content{
		Origin: test.FileOrigin{FilePath: filePath},
		Data:   data,
	}, nil
}
