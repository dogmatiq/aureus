package runner

import (
	"bytes"
	"io"
	"strings"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/test"
)

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	GenerateOutput OutputGenerator
}

// OutputGenerator produced output for a test case's input.
//
// in is the input data for the test case. out is meta-data about the expected
// output, which may be used to influence the kind of output that the generator
// writes to w.
type OutputGenerator func(
	w io.Writer,
	in test.Content,
	out test.ContentMetaData,
) error

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
	err := x.GenerateOutput(
		output,
		a.Input,
		a.Output.ContentMetaData,
	)

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
