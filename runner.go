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
	Output func(EqualAssertion) (string, error)
}

// Run makes the assertions described by all documents within a [TestSuite].
func (r *Runner[T]) Run(t T, s Node) {
	t.Helper()
	t.Run(
		s.Name,
		func(t T) {
			t.Helper()

			if s.Skip {
				t.SkipNow()
				return // return in case stubbed SkipNow() impementation does not panic
			}

			for _, c := range s.Children {
				r.Run(t, c)
			}

			for _, d := range s.Documents {
				r.RunDocument(t, d)
			}
		},
	)
}

// RunDocument makes the assertions described within a single [TestDocument].
func (r *Runner[T]) RunDocument(t T, d TestDocument) {
	t.Helper()
	t.Run(
		d.Name,
		func(t T) {
			t.Helper()

			if d.Skip {
				t.SkipNow()
				return // return in case stubbed SkipNow() impementation does not panic
			}

			for _, a := range d.Assertions {
				r.assert(t, a)
			}
		},
	)
}

func (r *Runner[T]) assert(t T, a EqualAssertion) {
	t.Helper()
	t.Run(
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

			output, err := r.Output(a)
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
