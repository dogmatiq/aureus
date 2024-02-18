package loader

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/dogmatiq/aureus/test"
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
func (b *TestBuilder) AddTest(t test.Test) {
	if len(t.SubTests) > 0 || t.Assertion != nil {
		b.tests = append(b.tests, t)
	}
}

// AddContent adds content to the builder.
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
