package fileloader

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/dogmatiq/aureus/internal/rootfs"
	"github.com/dogmatiq/aureus/loader"
	"github.com/dogmatiq/aureus/loader/internal/diriter"
	"github.com/dogmatiq/aureus/test"
)

// Loader loads [test.Test] values from files containing test inputs and
// expected outputs.
type Loader struct {
	options loadOptions
}

// NewLoader returns a new [Loader], which loads golden file tests from the
// filesystem.
func NewLoader(options ...LoadOption) *Loader {
	l := &Loader{
		options: loadOptions{
			FS:          rootfs.FS,
			Recurse:     true,
			LoadContent: LoadContent,
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
	return loadDir(opts, dir)
}

func loadDir(
	opts loadOptions,
	dirPath string,
) (test.Test, error) {
	var subTests []test.Test
	groups := map[string][]loader.ContentEnvelope{}

	err := diriter.Each(
		opts.FS,
		opts.Recurse,
		dirPath,
		func(dirPath string) error {
			t, err := loadDir(opts, dirPath)
			if err != nil {
				return err
			}
			if len(t.SubTests) != 0 {
				subTests = append(subTests, t)
			}
			return nil
		},
		func(filePath string) error {
			return loadFile(opts, filePath, groups)
		},
	)
	if err != nil {
		return test.Test{}, err
	}

	for n, envelopes := range groups {
		s, err := buildTest(n, envelopes)
		if err != nil {
			return test.Test{}, err
		}
		subTests = append(subTests, s)
	}

	slices.SortFunc(
		subTests,
		func(a, b test.Test) int {
			return strings.Compare(a.Name, b.Name)
		},
	)

	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	return test.New(
		name,
		test.WithSkip(skip),
		test.WithSubTests(subTests...),
	), nil
}

func loadFile(
	opts loadOptions,
	filePath string,
	groups map[string][]loader.ContentEnvelope,
) error {
	f, err := opts.FS.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	base := path.Base(filePath)
	name, skip := strings.CutPrefix(base, "_")

	c, err := opts.LoadContent(name, f)
	if err != nil {
		return err
	}

	groups[c.Group] = append(
		groups[c.Group],
		loader.ContentEnvelope{
			File:    filePath,
			Skip:    skip,
			Content: c,
		},
	)

	return nil
}

func buildTest(name string, envelopes []loader.ContentEnvelope) (test.Test, error) {
	inputs, outputs := loader.SeparateContentByRole(envelopes)

	switch {
	case len(inputs) == 0:
		return test.Test{}, fmt.Errorf("output file %q has no associated input files", outputs[0].File)
	case len(outputs) == 0:
		return test.Test{}, fmt.Errorf("input file %q has no associated output files", inputs[0].File)
	case len(inputs) == 1 && len(outputs) == 1:
		return buildSingleTest(name, inputs[0], outputs[0]), nil
	default:
		return buildMatrixTest(name, inputs, outputs), nil
	}
}

func buildSingleTest(name string, input, output loader.ContentEnvelope) test.Test {
	return test.New(
		name,
		test.WithSkip(input.Skip || output.Skip),
		test.WithAssertion(
			test.EqualAssertion{
				Input:  input.AsTestContent(),
				Output: output.AsTestContent(),
			},
		),
	)
}

func buildMatrixTest(name string, inputs, outputs []loader.ContentEnvelope) test.Test {
	t := test.New(name)

	testName := func(input, output loader.ContentEnvelope) string {
		if input.Content.Language != "" && output.Content.Language != "" {
			return fmt.Sprintf("%s -> %s", input.Content.Language, output.Content.Language)
		}

		return fmt.Sprintf(
			"%s -> %s",
			path.Base(input.File),
			path.Base(output.File),
		)
	}

	for _, output := range outputs {
		for _, input := range inputs {
			t.SubTests = append(
				t.SubTests,
				test.New(
					testName(input, output),
					test.WithSkip(input.Skip || output.Skip),
					test.WithAssertion(
						test.EqualAssertion{
							Input:  input.AsTestContent(),
							Output: output.AsTestContent(),
						},
					),
				),
			)
		}
	}

	return t
}
