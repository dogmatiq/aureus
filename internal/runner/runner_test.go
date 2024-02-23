package runner_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	. "github.com/dogmatiq/aureus/internal/runner"
	"github.com/dogmatiq/aureus/internal/test"
)

func TestRunner(t *testing.T) {
	loader := fileloader.NewLoader()

	// expected to pass
	{
		tst, err := loader.Load("testdata/pass")
		if err != nil {
			t.Fatal(err)
		}

		runner := &Runner[*testing.T]{
			GenerateOutput: func(
				t *testing.T,
				w io.Writer,
				in test.Content,
				out test.ContentMetaData,
			) {
				if err := prettyPrint(w, in, out); err != nil {
					t.Fatal(err)
				}
			},
		}

		runner.Run(t, tst)
	}

	// expected to fail
	{
		tst, err := loader.Load("testdata/fail")
		if err != nil {
			t.Fatal(err)
		}

		runner := &Runner[*testingT]{
			GenerateOutput: func(
				t *testingT,
				w io.Writer,
				in test.Content,
				out test.ContentMetaData,
			) {
				if err := prettyPrint(w, in, out); err != nil {
					t.Fatal(err)
				}
			},
		}

		x := &testingT{T: t}
		runner.Run(x, tst)

		for _, leaf := range x.leaves() {
			if !leaf.Failed {
				x.Errorf("expected %q to fail", leaf.Name())
			}
		}
	}
}

type testingT struct {
	*testing.T
	Children []*testingT
	Failed   bool
}

func (t *testingT) Run(name string, fn func(*testingT)) bool {
	return t.T.Run(
		name,
		func(x *testing.T) {
			child := &testingT{
				T: x,
			}

			t.Children = append(t.Children, child)

			fn(child)
		},
	)
}

func (t *testingT) Fail() {
	t.Failed = true
}

func (t *testingT) leaves() []*testingT {
	var leaves []*testingT

	if len(t.Children) == 0 {
		leaves = append(leaves, t)
	} else {
		for _, child := range t.Children {
			leaves = append(leaves, child.leaves()...)
		}
	}

	return leaves
}

func prettyPrint(
	w io.Writer,
	in test.Content,
	out test.ContentMetaData,
) error {
	if len(in.Data) == 0 {
		return nil
	}

	var v any
	if err := json.Unmarshal(in.Data, &v); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
