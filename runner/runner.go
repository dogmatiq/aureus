package runner

import (
	"strings"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/test"
)

// NativeRunner executes tests under Go's native testing framework.
type NativeRunner = Runner[*testing.T]

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	Output func(input test.Content) ([]byte, error)
}

// Run makes the assertions described by all documents within a [TestSuite].
func (r *Runner[T]) Run(t T, x test.Test) {
	t.Helper()
	t.Run(
		x.Name,
		func(t T) {
			t.Helper()

			if x.Flags.Has(test.FlagSkipped) {
				t.SkipNow()
				// Return in case stubbed SkipNow() impementation does not panic
				return
			}

			for _, s := range x.SubTests {
				r.Run(t, s)
			}

			if x.Assertion != nil {
				x.Assertion.AcceptVisitor(
					assertionExecutor[T]{r.Output, t},
					test.WithT(t),
				)
			}
		},
	)
}

// assertionExecutor is an impelmentation of [test.AssertionVisitor] that
// performs assertions within the context of a test.
type assertionExecutor[T TestingT[T]] struct {
	Output   func(input test.Content) ([]byte, error)
	TestingT T
}

func (x assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	var m strings.Builder

	m.WriteString("\n")
	m.WriteString("--- INPUT ---\n")
	m.Write(a.Input.Data)

	output, err := x.Output(a.Input)
	if err != nil {
		x.TestingT.Fail()
		m.WriteString("--- OUTPUT (error) ---\n")
		m.WriteString(err.Error())
	} else if d := diff.Diff(
		"want", a.Output.Data,
		"got", output,
	); d != nil {
		x.TestingT.Fail()
		m.WriteString("--- OUTPUT (-want +got) ---\n")
		m.WriteString(string(d))
	} else {
		m.WriteString("--- OUTPUT ---\n")
		m.Write(output)
	}

	x.TestingT.Log(m.String())
}
