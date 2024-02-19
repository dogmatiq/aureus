package runner_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/dogmatiq/aureus/loader/fileloader"
	. "github.com/dogmatiq/aureus/runner"
	"github.com/dogmatiq/aureus/test"
)

func TestRunner(t *testing.T) {
	loader := fileloader.NewLoader()

	formatJSON := func(input test.Content, output io.Writer) error {
		if len(input.Data) == 0 {
			return nil
		}

		var v any
		if err := json.Unmarshal(
			[]byte(input.Data),
			&v,
		); err != nil {
			return err
		}

		enc := json.NewEncoder(output)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	}

	// expected to pass
	{
		test, err := loader.Load("testdata/pass")
		if err != nil {
			t.Fatal(err)
		}

		runner := &Runner[*testing.T]{
			GenerateOutput: formatJSON,
		}

		runner.Run(t, test)
	}

	// expected to fail
	{
		test, err := loader.Load("testdata/fail")
		if err != nil {
			t.Fatal(err)
		}

		runner := &Runner[*testingT]{
			GenerateOutput: formatJSON,
		}

		x := &testingT{T: t}
		runner.Run(x, test)

		for _, leaf := range x.leaves() {
			if !leaf.Failed {
				t.Errorf("expected %q to fail", leaf.Name())
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
