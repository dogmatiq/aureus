package loader

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/dogmatiq/aureus/internal/natsort"
	"github.com/dogmatiq/aureus/internal/test"
)

// TestBuilder builds [test.Test] values from groups of correlated inputs and
// outputs.
type TestBuilder struct {
	tests  []test.Test
	groups map[string]*group
	anon   ContentEnvelope
}

type group struct {
	Name            string
	Inputs, Outputs []ContentEnvelope
}

// AddTest adds a pre-built test to the builder.
//
// Empty tests are ignored.
func (b *TestBuilder) AddTest(t test.Test) {
	if len(t.SubTests) > 0 || t.Assertion != nil {
		b.tests = append(b.tests, t)
	}
}

// AddContent adds content to the builder.
//
// Content with no [Role] is ignored.
func (b *TestBuilder) AddContent(env ContentEnvelope) error {
	if env.Content.Role == NoRole {
		return nil
	}

	if env.Content.Group == "" {
		return b.addAnonymousContent(env)
	}

	return b.addContent(env)
}

func (b *TestBuilder) addContent(env ContentEnvelope) error {
	switch b.anon.Content.Role {
	case Input:
		return NoOutputsError{[]ContentEnvelope{b.anon}}
	case Output:
		return NoInputsError{[]ContentEnvelope{b.anon}}
	}

	switch env.Content.Role {
	case Input:
		g := b.group(env.Content.Group)
		g.Inputs = append(g.Inputs, env)
	case Output:
		g := b.group(env.Content.Group)
		g.Outputs = append(g.Outputs, env)
	}

	return nil
}

func (b *TestBuilder) addAnonymousContent(env ContentEnvelope) error {
	emit := func(in, out ContentEnvelope) {
		name := fmt.Sprintf("anonymous test on line %d", out.Line)
		if out.Line == 0 {
			name = fmt.Sprintf("anonymous test in %s", path.Base(out.File))
		}

		g := b.group(name)
		g.Inputs = append(g.Inputs, in)
		g.Outputs = append(g.Outputs, out)
		b.anon = ContentEnvelope{}
	}

	switch b.anon.Content.Role {
	case NoRole:
		b.anon = env
	case Input:
		if env.Content.Role == Input {
			return NoOutputsError{[]ContentEnvelope{b.anon}}
		}
		emit(b.anon, env)
	case Output:
		if env.Content.Role == Output {
			return NoInputsError{[]ContentEnvelope{b.anon}}
		}
		emit(env, b.anon)
	}

	return nil
}

// Build returns tests built from the inputs and outputs, sorted by name.
func (b *TestBuilder) Build() ([]test.Test, error) {
	tests := make([]test.Test, 0, len(b.groups)+len(b.tests))
	tests = append(tests, b.tests...)

	for _, g := range b.groups {
		t, err := buildTest(g)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}

	slices.SortFunc(
		tests,
		func(a, b test.Test) int {
			if natsort.Less(a.Name, b.Name) {
				return -1
			} else if natsort.Less(b.Name, a.Name) {
				return 1
			} else {
				return 0
			}
		},
	)

	return tests, nil
}

// LoadDir loads tests from the given directory.
func LoadDir(
	fsys fs.FS,
	dirPath string,
	recurse bool,
	build func(*TestBuilder, fs.FS, string) error,
) (test.Test, error) {
	var builder TestBuilder

	entries, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return test.Test{}, err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := path.Join(dirPath, entry.Name())

		if entry.IsDir() {
			if recurse {
				t, err := LoadDir(fsys, entryPath, true, build)
				if err != nil {
					return test.Test{}, err
				}
				builder.AddTest(t)
			}
		} else {
			if err := build(&builder, fsys, entryPath); err != nil {
				return test.Test{}, err
			}
		}
	}

	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	subTests, err := builder.Build()
	if err != nil {
		return test.Test{}, err
	}

	return test.New(
		name,
		test.WithSkip(skip),
		test.WithSubTests(subTests...),
	), nil
}

// group returns the group with the given name, creating it if necessary.
func (b *TestBuilder) group(name string) *group {
	if b.groups == nil {
		b.groups = map[string]*group{}
	} else if g, ok := b.groups[name]; ok {
		return g
	}

	g := &group{Name: name}
	b.groups[name] = g
	return g
}

// buildTest builds a test for the given group.
func buildTest(g *group) (test.Test, error) {
	inputs := len(g.Inputs)
	outputs := len(g.Outputs)

	if inputs == 0 {
		return test.Test{}, NoInputsError{g.Outputs}
	} else if outputs == 0 {
		return test.Test{}, NoOutputsError{g.Inputs}
	} else if inputs == 1 && outputs == 1 {
		return buildSingleTest(g), nil
	} else {
		return buildMatrixTest(g), nil
	}
}

func buildSingleTest(g *group) test.Test {
	input := g.Inputs[0]
	output := g.Outputs[0]

	return test.New(
		g.Name,
		test.WithSkip(input.Skip || output.Skip),
		test.WithAssertion(
			test.EqualAssertion{
				Input:  input.AsTestContent(),
				Output: output.AsTestContent(),
			},
		),
	)
}

func buildMatrixTest(g *group) test.Test {
	t := test.New(g.Name)

	testName := func(input, output ContentEnvelope) string {
		if input.Content.Language != "" && output.Content.Language != "" {
			return fmt.Sprintf("%s -> %s", input.Content.Language, output.Content.Language)
		}

		return fmt.Sprintf(
			"%s -> %s",
			path.Base(input.File),
			path.Base(output.File),
		)
	}

	for _, output := range g.Outputs {
		for _, input := range g.Inputs {
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
