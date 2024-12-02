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
	GenerateOutput  OutputGenerator[T]
	TrimSpace       bool // TODO: make this a loader concern
	BlessStrategy   BlessStrategy
	AssertionFilter func(test.Assertion) bool
	PackagePath     string
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
		fmt.Sprintf("INPUT (%s)", location(a.Input)),
		a.Input.Data,
		"\x1b[2m",
	)

	if r.AssertionFilter != nil && !r.AssertionFilter(a) {
		t.SkipNow()
		return
	}

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

	want := a.Output.Data
	got, err := io.ReadAll(f)
	if err != nil {
		t.Log("unable to read output file:", err)
	}

	if r.TrimSpace {
		want = append(bytes.TrimRight(want, "\n"), '\n')
		got = append(bytes.TrimRight(got, "\n"), '\n')
	}

	diff := diff.ColorDiff(
		location(a.Output),
		want,
		f.Name(),
		got,
	)

	messages := []string{
		"\x1b[1mTo run this test again, use:\n\n" +
			"    \x1b[2m" + r.goTestCommand(t) + "\x1b[0m",
	}

	if len(diff) == 0 {
		logSection(
			t,
			fmt.Sprintf("OUTPUT (%s)", location(a.Output)),
			a.Output.Data,
			"\x1b[33;2m",
			messages...,
		)
		return
	}

	switch r.BlessStrategy {
	case BlessAvailable:
		t.Fail()
		messages = append(
			messages,
			"\x1b[1mTo \x1b[33maccept this output as correct\x1b[37m from now on, add the \x1b[33m-aureus.bless\x1b[37m flag:\n\n"+
				"    \x1b[2m"+r.goTestCommand(t)+" -aureus.bless\x1b[0m",
		)

	case BlessDisabled:
		t.Fail()

	case BlessEnabled:
		if err := bless(a.Output, got); err != nil {
			t.Log("unable to bless output:", err)
			t.Fail()
			return
		}

		messages = append(
			messages,
			"\x1b[1mThe current \x1b[33moutput has been blessed\x1b[0m. Future runs will consider this output correct.\x1b[0m",
		)
	}

	logSection(
		t,
		"OUTPUT DIFF",
		diff,
		"",
		messages...,
	)

}

func location(c test.Content) string {
	if c.IsEntireFile() {
		return c.File
	}
	return fmt.Sprintf("%s:%d", c.File, c.Line)
}

func log(t LoggerT, fn func(w *strings.Builder)) {
	t.Helper()
	var w strings.Builder
	fn(&w)
	t.Log(w.String())
}

func logSection(
	t LoggerT,
	title string,
	body []byte,
	bodyANSI string,
	messages ...string,
) {
	t.Helper()

	log(t, func(w *strings.Builder) {
		w.WriteString("\x1b[0m")

		w.WriteString("\n")
		w.WriteString("\n")

		w.WriteString("\x1b[1m")
		w.WriteString("╭────")

		w.WriteString("\x1b[7m") // inverse
		w.WriteString(" ")
		w.WriteString(title)

		w.WriteString(" ")
		w.WriteString("\x1b[27m") // reset inverse
		w.WriteString("────\x1b[0m────\x1b[2m──┈\x1b[0m\n")

		w.WriteString("\x1b[1m│\x1b[0m\n")

		for _, line := range bytes.Split(body, newLine) {
			w.WriteString("\x1b[1m│\x1b[0m  ")
			w.WriteString(bodyANSI)
			w.Write(line)
			w.WriteString("\x1b[0m\n")
		}

		w.WriteString("\x1b[1m│\x1b[0m\n")
		w.WriteString("\x1b[1m╰────\x1b[0m────\x1b[2m──┈\x1b[0m\n")

		for _, t := range messages {
			w.WriteString("\n")
			w.WriteString("\x1b[33m✦\x1b[0m ")
			w.WriteString(t)
			w.WriteString("\x1b[0m\n")
		}
	})
}

func (r *Runner[T]) computeDiff(
	want test.Content,
	got *os.File,
) ([]byte, error) {
	data, err := io.ReadAll(got)
	if err != nil {
		return nil, fmt.Errorf("unable to read output file: %w", err)
	}

	if r.TrimSpace {
		data = append(bytes.TrimSpace(data), '\n')
	}

	return diff.ColorDiff(
		location(want),
		want.Data,
		got.Name(),
		data,
	), nil
}

var (
	separator = strings.Repeat("=", 10)
	newLine   = []byte("\n")
)
