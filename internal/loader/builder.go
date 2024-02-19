package loader

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/dogmatiq/aureus/internal/test"
)

// TestBuilder builds [test.Test] values from groups of correlated inputs and
// outputs.
type TestBuilder struct {
	tests  []test.Test
	groups map[string]*group
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
func (b *TestBuilder) AddContent(env ContentEnvelope) {
	switch env.Content.Role {
	case Input:
		g := b.group(env.Content.Group)
		g.Inputs = append(g.Inputs, env)
	case Output:
		g := b.group(env.Content.Group)
		g.Outputs = append(g.Outputs, env)
	}
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
			return strings.Compare(a.Name, b.Name)
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

// location returns a string that describes the location of the given content.
func location(env ContentEnvelope) string {
	if env.Line > 0 {
		return fmt.Sprintf("%s:%d", env.File, env.Line)
	}
	return env.File
}

// buildTest builds a test for the given group.
func buildTest(g *group) (test.Test, error) {
	switch {
	case len(g.Inputs) == 0:
		return test.Test{}, fmt.Errorf("output from %s has no associated inputs", location(g.Outputs[0]))
	case len(g.Outputs) == 0:
		return test.Test{}, fmt.Errorf("input from %s has no associated outputs", location(g.Inputs[0]))
	case len(g.Inputs) == 1 && len(g.Outputs) == 1:
		return buildSingleTest(g), nil
	default:
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
