package runner

import (
	"bytes"
	"io"
	"strings"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/test"
)

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	GenerateOutput OutputGenerator
}

// OutputGenerator is a function that generates output for a given input.
type OutputGenerator func(input test.Content, output io.Writer) error

// Run makes the assertions described by all documents within a [TestSuite].
func (r *Runner[T]) Run(t T, x test.Test) {
	t.Helper()
	t.Run(
		x.Name,
		func(t T) {
			t.Helper()

			if x.Skip {
				t.SkipNow()
				// Return in case stubbed SkipNow() impementation does not panic
				return
			}

			for _, s := range x.SubTests {
				r.Run(t, s)
			}

			if x.Assertion != nil {
				x.Assertion.AcceptVisitor(
					assertionExecutor[T]{r.GenerateOutput, t},
					test.WithT(t),
				)
			}
		},
	)
}

// assertionExecutor is an impelmentation of [test.AssertionVisitor] that
// performs assertions within the context of a test.
type assertionExecutor[T TestingT[T]] struct {
	GenerateOutput OutputGenerator
	TestingT       T
}

func (x assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	var m strings.Builder

	m.WriteString("\n")
	m.WriteString("--- INPUT ---\n")
	m.WriteString(a.Input.Data)

	output := &bytes.Buffer{}
	err := x.GenerateOutput(a.Input, output)

	if err != nil {
		x.TestingT.Fail()
		m.WriteString("--- OUTPUT (error) ---\n")
		m.WriteString(err.Error())
	} else if d := diff.Diff(
		"want", []byte(a.Output.Data),
		"got", output.Bytes(),
	); d != nil {
		x.TestingT.Fail()
		m.WriteString("--- OUTPUT (-want +got) ---\n")
		m.Write(d)
	} else {
		m.WriteString("--- OUTPUT ---\n")
		output.WriteTo(&m)
	}

	x.TestingT.Log(m.String())
}
