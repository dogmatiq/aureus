package runner

import (
	"fmt"
	"io"

	"github.com/dogmatiq/aureus/internal/test"
)

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	GenerateOutput OutputGenerator[T]
	TrimSpace      bool
}

// OutputGenerator produced output for a test case's input.
//
// in is the input data for the test case. out is meta-data about the expected
// output, which may be used to influence the kind of output that the generator
// writes to w.
type OutputGenerator[T any] func(
	t T,
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
					&assertionExecutor[T]{t, r},
					test.WithT(t),
				)
			}
		},
	)
}

func location(c test.Content) string {
	if c.Line == 0 {
		return c.File
	}
	return fmt.Sprintf("%s:%d", c.File, c.Line)
}
