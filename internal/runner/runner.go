package runner

import (
	"fmt"

	"github.com/dogmatiq/aureus/internal/test"
)

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	GenerateOutput OutputGenerator[T]
	TrimSpace      bool
}

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

			e := assertionExecutor[T]{t, r}
			for _, a := range x.Assertions {
				a.AcceptVisitor(&e, test.WithT(t))
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
