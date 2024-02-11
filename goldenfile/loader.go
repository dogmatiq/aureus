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
				correlate(in, out)
			}
		} else {
			in := &inputFile{Path: n}
			inputs = append(inputs, in)
			for _, out := range outputs {
				correlate(in, out)
			}
		}
	}

	for _, out := range outputs {
		if len(out.Inputs) == 0 {
			return test.Test{}, fmt.Errorf("output file %q has no associated input files", out.Path)
		}

		if out.IsMatrix {
			panic("not implemented")
		}

		child, err := buildSingleTest(opts, out, inherited)
		if err != nil {
			return test.Test{}, err
		}

		parent.SubTests = append(parent.SubTests, child)
	}

	for _, in := range inputs {
		if len(in.Outputs) == 0 {
			return test.Test{}, fmt.Errorf("input file %q has no associated output files", in.Path)
		}
	}

	return parent, nil
}

func buildSingleTest(
	opts loadOptions,
	out *outputFile,
	inherited test.FlagSet,
) (test.Test, error) {
	if len(out.Inputs) != 1 {
		panic("unexpected number of inputs")
	}

	output, err := loadContent(opts, out.Path)
	if err != nil {
		return test.Test{}, err
	}

	in := out.Inputs[0]
	input, err := loadContent(opts, in.Path)
	if err != nil {
		return test.Test{}, err
	}

	t, _ := test.New(
		test.WithNameFromPath(out.Path),
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

// func loadInput(
// 	opts loadOptions,
// 	filePath string,
// 	output test.Content,
// 	inherited test.FlagSet,
// ) (test.Test, error) {

// 	t, _ := test.New(
// 		test.WithNameFromPath(filePath),
// 		test.WithInheritedFlags(inherited),
// 		test.WithAssertion(
// 			test.EqualAssertion{
// 				Input:  inputFile,
// 				Output: output,
// 			},
// 		),
// 	)

// 	return t, nil
// }

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
		Data: string(data),
		File: filePath,
	}, nil
}

type outputFile struct {
	Path     string
	IsInput  InputPredicate
	IsMatrix bool
	Inputs   []*inputFile
}

type inputFile struct {
	Path     string
	IsMatrix bool
	Outputs  []*outputFile
}

func correlate(in *inputFile, out *outputFile) {
	if out.IsInput(path.Base(in.Path)) {
		in.Outputs = append(in.Outputs, out)
		out.Inputs = append(out.Inputs, in)

		if len(in.Outputs) > 1 || len(out.Inputs) > 1 {
			in.IsMatrix = true
			out.IsMatrix = true
		}
	}
}
