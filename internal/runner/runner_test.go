package runner_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	. "github.com/dogmatiq/aureus/internal/runner"
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
			GenerateOutput: prettyPrint[*testing.T],
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
			GenerateOutput: prettyPrint[*testingT],
		}

		x := &testingT{T: t}
		runner.Run(x, tst)

		for _, leaf := range x.leaves() {
			if !leaf.Failed() {
				x.Errorf("expected %q to fail", leaf.Name())
			}
		}
	}
}

type testingT struct {
	*testing.T
	Children []*testingT
	failed   bool
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
	t.failed = true
}

func (t *testingT) Failed() bool {
	return t.failed
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

func prettyPrint[T any](_ T, in Input, out Output) error {
	var v any
	dec := json.NewDecoder(in)
	if err := dec.Decode(&v); err != nil {
		if err == io.EOF { // handle "empty" test
			return nil
		}
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
