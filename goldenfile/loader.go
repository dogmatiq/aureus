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
	parent, inherited := test.New(
		test.WithNameFromPath(dirPath),
		test.WithInheritedFlags(inherited),
	)

	entries, err := fs.ReadDir(opts.FS, dirPath)
	if err != nil {
		return test.Test{}, err
	}

	type output struct {
		FilePath string
		IsInput  InputPredicate
	}

	type input struct {
		FilePath        string
		MatchedToOutput bool
	}

	var (
		inputs  []*input
		outputs []*output
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
			outputs = append(outputs, &output{n, isInput})
		} else {
			inputs = append(inputs, &input{n, false})
		}
	}

	for _, out := range outputs {
		var matching []string
		for _, in := range inputs {
			if out.IsInput(path.Base(in.FilePath)) {
				in.MatchedToOutput = true
				matching = append(matching, in.FilePath)
			}
		}

		if len(matching) == 0 {
			return test.Test{}, fmt.Errorf("output file %q has no associated input files", out.FilePath)
		}

		child, err := loadOutput(opts, out.FilePath, matching, inherited)
		if err != nil {
			return test.Test{}, err
		}

		parent.SubTests = append(parent.SubTests, child)
	}

	for _, in := range inputs {
		if !in.MatchedToOutput {
			return test.Test{}, fmt.Errorf("input file %q has no associated output files", in.FilePath)
		}
	}

	return parent, nil
}

func loadOutput(
	opts loadOptions,
	filePath string,
	inputs []string,
	inherited test.FlagSet,
) (test.Test, error) {
	parent, inherited := test.New(
		test.WithNameFromPath(filePath),
		test.WithInheritedFlags(inherited),
	)

	output, err := loadContent(opts, filePath)
	if err != nil {
		return test.Test{}, err
	}

	for _, filePath := range inputs {
		child, err := loadInput(opts, filePath, output, inherited)
		if err != nil {
			return test.Test{}, err
		}
		parent.SubTests = append(parent.SubTests, child)
	}

	return parent, nil
}

func loadInput(
	opts loadOptions,
	filePath string,
	output test.Content,
	inherited test.FlagSet,
) (test.Test, error) {
	input, err := loadContent(opts, filePath)
	if err != nil {
		return test.Test{}, err
	}

	t, _ := test.New(
		test.WithNameFromPath(filePath),
		test.WithInheritedFlags(inherited),
		test.WithAssertion(
			test.EqualAssertion{
				Input:  input,
				Output: output,
			},
		),
	)

	return t, nil
}

// loadContent loads the content of an input file or output file.
func loadContent(
	opts loadOptions,
	filePath string,
) (test.Content, error) {
	data, err := fs.ReadFile(opts.FS, filePath)
	if err != nil {
		return test.Content{}, err
	}

	return test.Content{
		Data: data,
		File: filePath,
	}, nil
}
