package runner

import (
	"fmt"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/dogmatiq/aureus/test"
)

// NativeRunner executes tests under Go's native testing framework.
type NativeRunner = Runner[*testing.T]

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	Output func(input test.Content) (string, error)
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

			for _, a := range x.Assertions {
				// NOTE: We don't use [test.AssertionVisitor] here because it
				// adds an extra method call (AcceptVisitor) that we're unable
				// to mark as a helper with t.Helper().
				switch a := a.(type) {
				case *test.EqualAssertion:
					r.equal(t, a)
				// case *test.DiffAssertion:
				// 	r.diff(t, a)
				default:
					panic(fmt.Sprintf("unsupported assertion type (%T)", a))
				}
			}
		},
	)
}

func (r *Runner[T]) equal(t T, a *test.EqualAssertion) {
	t.Helper()
	var m strings.Builder

	m.WriteString("\n")
	m.WriteString("--- INPUT ---\n")
	m.WriteString(a.Input.Data)

	output, err := r.Output(a.Input)
	if err != nil {
		t.Fail()
		m.WriteString("--- OUTPUT (error) ---\n")
		m.WriteString(err.Error())
	} else if output != a.Output.Data {
		t.Fail()
		m.WriteString("--- OUTPUT (-want +got) ---\n")
		m.WriteString(diff.LineDiff(a.Output.Data, output))
	} else {
		m.WriteString("--- OUTPUT ---\n")
		m.WriteString(output)
	}

	t.Log(m.String())
}
