package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/dogmatiq/aureus/internal/diff"
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
)

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
					assertionExecutor[T]{
						t,
						r.GenerateOutput,
						r.TrimSpace,
					},
					test.WithT(t),
				)
			}
		},
	)
}

// assertionExecutor is an impelmentation of [test.AssertionVisitor] that
// performs assertions within the context of a test.
type assertionExecutor[T TestingT[T]] struct {
	TestingT       T
	GenerateOutput OutputGenerator[T]
	TrimSpace      bool
}

func (x assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	x.log("=== BEGIN INPUT (%s) ===", location(a.Input))
	x.log("%s", string(x.sanitizeForLog(a.Input.Data)))
	x.log("=== END INPUT ===")

	output := &bytes.Buffer{}
	x.GenerateOutput(
		x.TestingT,
		output,
		a.Input,
		a.Output.ContentMetaData,
	)

	f, err := os.CreateTemp("", "aureus-")
	if err != nil {
		x.TestingT.Log("failed to create temporary file:", err)
		x.TestingT.Fail()
		return
	}
	defer f.Close()

	if d := diff.Diff(
		location(a.Output), x.sanitizeForDiff(a.Output.Data),
		f.Name(), x.sanitizeForDiff(output.Bytes()),
	); d != nil {
		x.TestingT.Fail()
		x.log("=== BEGIN OUTPUT DIFF (-want +got) ===")
		x.log("%s", d)
		x.log("=== END OUTPUT DIFF ===")

		if _, err := f.Write(a.Output.Data); err != nil {
			x.TestingT.Log("unable to write output to temporary file:", err)
			os.Remove(f.Name())
			return
		}
	} else {
		os.Remove(f.Name())
		x.log("=== BEGIN OUTPUT (%s) ===", location(a.Output))
		x.log("%s", x.sanitizeForLog(a.Output.Data))
		x.log("=== END OUTPUT ===")
	}
}

func (x assertionExecutor[T]) log(format string, args ...any) {
	x.TestingT.Helper()

	lines := strings.Split(
		fmt.Sprintf(format, args...),
		"\n",
	)

	for _, l := range lines {
		x.TestingT.Log(l)
	}
}

func (x assertionExecutor[T]) sanitizeForLog(data []byte) []byte {
	if x.TrimSpace {
		return bytes.TrimSpace(data)
	}
	return data
}

func (x assertionExecutor[T]) sanitizeForDiff(data []byte) []byte {
	if x.TrimSpace {
		data = slices.Clone(data)
		data = bytes.TrimSpace(data)
		return append(data, '\n')
	}
	return slices.Clone(data)
}

func location(c test.Content) string {
	if c.Line == 0 {
		return c.File
	}
	return fmt.Sprintf("%s:%d", c.File, c.Line)
}
