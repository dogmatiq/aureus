package goldenfile

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
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
			LoadFile: LoadFile,
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
	return l.loadDir(dir, test.EmptyFlagSet)
}

func (l *Loader) loadDir(
	dirPath string,
	flags test.FlagSet,
) (test.Test, error) {
	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	t, flags := test.New(
		test.WithName(name),
		test.If(skip, test.WithFlag(test.FlagSkipped)),
		test.WithInheritedFlags(flags),
	)

	entries, err := fs.ReadDir(l.options.FS, dirPath)
	if err != nil {
		return test.Test{}, err
	}

	filesByTest := map[string][]File{}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := path.Join(dirPath, entry.Name())

		if entry.IsDir() {
			if l.options.Recurse {
				s, err := l.loadDir(entryPath, flags)
				if err != nil {
					return test.Test{}, err
				}
				t.SubTests = append(t.SubTests, s)
			}
		} else {
			file, ok, err := l.options.LoadFile(l.options.FS, entryPath)
			if err != nil {
				return test.Test{}, err
			}
			if ok {
				filesByTest[file.TestName] = append(filesByTest[file.TestName], file)
			}
		}
	}

	// Build a sub-test for each separate group of files.
	for n, files := range filesByTest {
		s, err := l.buildTest(n, files, flags)
		if err != nil {
			return test.Test{}, err
		}
		t.SubTests = append(t.SubTests, s)
	}

	// Sort by name those sub-tests that were built from file groups. This
	// ensures that the order of the sub-tests is deterministic, and also that
	// sub-tests build from directories appear first.
	slices.SortFunc(
		t.SubTests[len(t.SubTests)-len(filesByTest):],
		func(a, b test.Test) int {
			return strings.Compare(a.Name, b.Name)
		},
	)

	return t, nil
}

func (l *Loader) buildTest(
	name string,
	files []File,
	flags test.FlagSet,
) (test.Test, error) {
	var inputs, outputs []File
	for _, f := range files {
		if f.IsInput {
			inputs = append(inputs, f)
		} else {
			outputs = append(outputs, f)
		}
	}

	switch {
	case len(inputs) == 0:
		return test.Test{}, fmt.Errorf("output file %q has no associated input files", outputs[0].Content.File)
	case len(outputs) == 0:
		return test.Test{}, fmt.Errorf("input file %q has no associated output files", inputs[0].Content.File)
	case len(inputs) == 1 && len(outputs) == 1:
		return l.buildSingleTest(name, inputs[0], outputs[0], flags)
	default:
		return l.buildMatrixTest(name, inputs, outputs, flags)
	}
}

func (l *Loader) buildSingleTest(
	name string,
	input, output File,
	flags test.FlagSet,
) (test.Test, error) {
	t, _ := test.New(
		test.WithName(name),
		test.WithInheritedFlags(flags),
		test.If(
			input.IsSkipped || output.IsSkipped,
			test.WithFlag(test.FlagSkipped),
		),
		test.WithAssertion(
			test.EqualAssertion{
				Input:  input.Content,
				Output: output.Content,
			},
		),
	)

	return t, nil
}

func (l *Loader) buildMatrixTest(
	name string,
	inputs, outputs []File,
	flags test.FlagSet,
) (test.Test, error) {
	parent, _ := test.New(
		test.WithName(name),
		test.WithInheritedFlags(flags),
	)

	testName := func(input, output File) string {
		if input.Content.Language != "" && output.Content.Language != "" {
			return fmt.Sprintf("%s -> %s", input.Content.Language, output.Content.Language)
		}

		return fmt.Sprintf(
			"%s -> %s",
			path.Base(input.Content.File),
			path.Base(output.Content.File),
		)
	}

	for _, output := range outputs {
		for _, input := range inputs {
			t, _ := test.New(
				test.WithName(testName(input, output)),
				test.WithInheritedFlags(flags),
				test.If(
					input.IsSkipped || output.IsSkipped,
					test.WithFlag(test.FlagSkipped),
				),
				test.WithAssertion(
					test.EqualAssertion{
						Input:  input.Content,
						Output: output.Content,
					},
				),
			)
			parent.SubTests = append(parent.SubTests, t)
		}
	}

	return parent, nil
}
