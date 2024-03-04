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

var (
	separator = strings.Repeat("=", 10)
	newLine   = []byte("\n")
)

// assertionExecutor is an impelmentation of [test.AssertionVisitor] that
// performs assertions within the context of a test.
type assertionExecutor[T TestingT[T]] struct {
	TestingT T
	Runner   *Runner[T]
}

func (x *assertionExecutor[T]) generateOutput(in, out test.Content) (_ *os.File, err error) {
	f, err := os.CreateTemp("", "aureus-")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary file: %w", err)
	}
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	if err := x.Runner.GenerateOutput(
		x.TestingT,
		&input{
			Reader: bytes.NewReader(in.Data),
			meta:   in.ContentMetaData,
		},
		&output{
			Writer: f,
			meta:   out.ContentMetaData,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to generate output: %w", err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("unable to seek to beginning of output file: %w", err)
	}

	return f, nil
}

func (x *assertionExecutor[T]) log(fn func(w *strings.Builder)) {
	x.TestingT.Helper()

	var w strings.Builder
	fn(&w)
	x.TestingT.Log(w.String())
}

func (x *assertionExecutor[T]) logSection(name string, data []byte, extra ...any) {
	x.TestingT.Helper()

	x.log(func(w *strings.Builder) {
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

	if x.Runner.TrimSpace {
		data = bytes.TrimSpace(data)
	}
	for _, line := range bytes.Split(data, newLine) {
		x.TestingT.Log(string(line))
	}

	x.log(func(w *strings.Builder) {
		w.WriteString(separator)
		w.WriteString(" END ")
		w.WriteString(name)
		w.WriteString(separator)
	})
}

func (x *assertionExecutor[T]) computeDiff(
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

	if x.Runner.TrimSpace {
		wantData = append(wantData, '\n')
		gotData = append(gotData, '\n')
	}

	return diff.Diff(wantLoc, wantData, gotLoc, gotData), nil
}
