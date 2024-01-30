package aureus

import (
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

// TestingT is a constraint for a type that is compatible with [testing.T].
type TestingT[T any] interface {
	Helper()
	Parallel()
	Run(string, func(T)) bool
	Log(...any)
	SkipNow()
	Fail()
}

// NativeRunner runs Aureus tests under Go's native testing framework.
type NativeRunner = Runner[*testing.T]

// Runner executes Aureus tests under any test framework with an interface
// similar to Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	Output func(Assertion) (string, error)
}

// Run executes performs the assertions within the given [Test] using t.
func (r *Runner[T]) Run(t T, test Test) {
	t.Helper()
	test.AcceptVisitor(&runner[T]{
		Runner:   r,
		TestingT: t,
	})
}

type runner[T TestingT[T]] struct {
	Runner   *Runner[T]
	TestingT T
}

func (r *runner[T]) VisitSuite(s Suite) {
	r.TestingT.Helper()
	r.TestingT.Run(
		s.Name,
		func(t T) {
			t.Helper()

			if s.Skip {
				t.SkipNow()
				return // return in case stubbed SkipNow() impementation does not panic
			}

			for _, sub := range s.Tests {
				r.Runner.Run(t, sub)
			}
		},
	)
}

func (r *runner[T]) VisitDocument(d Document) {
	r.TestingT.Helper()
	r.TestingT.Run(
		d.Name,
		func(t T) {
			t.Helper()

			if d.Skip {
				t.SkipNow()
				return // return in case stubbed SkipNow() impementation does not panic
			}

			for _, a := range d.Assertions {
				r.Runner.Run(t, a)
			}
		},
	)
}

func (r *runner[T]) VisitAssertion(a Assertion) {
	r.TestingT.Helper()
	r.TestingT.Run(
		a.Name,
		func(t T) {
			t.Helper()

			if a.Skip {
				t.SkipNow()
				return // return in case stubbed SkipNow() impementation does not panic
			}

			var m strings.Builder

			m.WriteString("\n")
			m.WriteString("--- INPUT (")
			m.WriteString(a.InputLanguage)
			m.WriteString(") ---\n")
			m.WriteString(a.Input)

			output, err := r.Runner.Output(a)
			if err != nil {
				t.Fail()
				m.WriteString("--- OUTPUT (error) ---\n")
				m.WriteString(err.Error())
			} else if output != a.ExpectedOutput {
				t.Fail()
				m.WriteString("--- OUTPUT (")
				m.WriteString(a.OutputLanguage)
				m.WriteString(" diff, -want +got) ---\n")
				m.WriteString(diff.LineDiff(a.ExpectedOutput, output))
			} else {
				m.WriteString("--- OUTPUT (")
				m.WriteString(a.OutputLanguage)
				m.WriteString(") ---\n")
				m.WriteString(output)
			}

			t.Log(m.String())
		},
	)
}
