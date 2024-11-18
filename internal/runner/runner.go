package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/test"
)

// Runner executes tests under any test framework with an interface similar to
// Go's native [*testing.T].
type Runner[T TestingT[T]] struct {
	GenerateOutput OutputGenerator[T]
	TrimSpace      bool // TODO: make this a loader concern
	BlessStrategy  BlessStrategy
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
				// TODO: this is here because the stubbed SkipNow()
				// impementation does not panic, can we make it unnecessary?
				return
			}

			for _, s := range x.SubTests {
				r.Run(t, s)
			}

			for _, a := range x.Assertions {
				r.assert(t, a)
			}
		},
	)
}

func (r *Runner[T]) assert(t T, a test.Assertion) {
	t.Helper()
	logSection(
		t,
		"INPUT",
		a.Input.Data,
		r.TrimSpace,
		location(a.Input),
	)

	f, err := generateOutput(t, r.GenerateOutput, a.Input, a.Output)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	defer func() {
		f.Close()
		if !t.Failed() {
			os.Remove(f.Name())
		}
	}()

	diff, err := computeDiff(
		r.TrimSpace,
		location(a.Output), bytes.NewReader(a.Output.Data),
		f.Name(), f,
	)
	if err != nil {
		t.Log("unable to generate diff:", err)
		t.Fail()
		return
	}

	if diff == nil {
		logSection(t, "OUTPUT", a.Output.Data, r.TrimSpace, location(a.Output))
		return
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Log("unable to rewind output file:", err)
		t.Fail()
		return
	}

	logSection(t, "OUTPUT DIFF", diff, r.TrimSpace)
	r.BlessStrategy.bless(t, a, f)
	t.Fail()
}

func location(c test.Content) string {
	if c.IsEntireFile() {
		return c.File
	}
	return fmt.Sprintf("%s:%d", c.File, c.Line)
}

func log[T TestingT[T]](t T, fn func(w *strings.Builder)) {
	t.Helper()
	var w strings.Builder
	fn(&w)
	t.Log(w.String())
}

func logSection[T TestingT[T]](
	t T,
	name string,
	data []byte,
	trimSpace bool,
	extra ...any,
) {
	t.Helper()

	log(t, func(w *strings.Builder) {
		w.WriteString(separator)
		w.WriteString(" BEGIN ")
		w.WriteString(name)

		if len(extra) > 0 {
			w.WriteString(" (")
			for _, v := range extra {
				fmt.Fprint(w, v)
			}
			w.WriteByte(')')
		}

		w.WriteByte(' ')
		w.WriteString(separator)
	})

	if trimSpace {
		data = bytes.TrimSpace(data)
	}
	for _, line := range bytes.Split(data, newLine) {
		t.Log(string(line))
	}

	log(t, func(w *strings.Builder) {
		w.WriteString(separator)
		w.WriteString(" END ")
		w.WriteString(name)
		w.WriteByte(' ')
		w.WriteString(separator)
		w.WriteByte('\n')
	})
}

func computeDiff(
	trimSpace bool,
	wantLoc string, want io.Reader,
	gotLoc string, got io.Reader,
) ([]byte, error) {
	wantData, err := io.ReadAll(want)
	if err != nil {
		return nil, fmt.Errorf("unable to read expected output: %w", err)
	}

	gotData, err := io.ReadAll(got)
	if err != nil {
		return nil, fmt.Errorf("unable to read output: %w", err)
	}

	if trimSpace {
		wantData = append(wantData, '\n')
		gotData = append(gotData, '\n')
	}

	return diff.ColorDiff(wantLoc, wantData, gotLoc, gotData), nil
}
