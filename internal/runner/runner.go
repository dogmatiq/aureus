package runner

import (
	"bytes"
	"fmt"
	"io"
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
						r.GenerateOutput,
						r.TrimSpace,
						t,
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
	GenerateOutput OutputGenerator[T]
	TrimSpace      bool
	TestingT       T
}

func (x assertionExecutor[T]) sanitize(data []byte) []byte {
	if x.TrimSpace {
		data = slices.Clone(data)
		return append(bytes.TrimSpace(data), '\n')
	}
	return data
}

func (x assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	var m strings.Builder

	input := x.sanitize(a.Input.Data)

	m.WriteString("=== INPUT (")
	m.WriteString(location(a.Input))
	m.WriteString(") ===\n")
	m.Write(input)
	x.TestingT.Log(m.String())
	m.Reset()

	output := &bytes.Buffer{}
	x.GenerateOutput(
		x.TestingT,
		output,
		a.Input,
		a.Output.ContentMetaData,
	)

	want := x.sanitize(a.Output.Data)
	got := x.sanitize(output.Bytes())

	if d := diff.Diff(
		location(a.Output), want,
		"generated-output", got,
	); d != nil {
		x.TestingT.Fail()
		m.WriteString("=== OUTPUT (-want +got) ===\n")
		m.Write(d)
	} else {
		m.WriteString("=== OUTPUT (")
		m.WriteString(location(a.Input))
		m.WriteString(" ===\n")
		output.WriteTo(&m)
	}

	x.TestingT.Log(m.String())
}

func location(c test.Content) string {
	if c.Line == 0 {
		return c.File
	}
	return fmt.Sprintf("%s:%d", c.File, c.Line)
}
