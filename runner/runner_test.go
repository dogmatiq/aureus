package runner_test

import (
	"encoding/json"
	"testing"

	"github.com/dogmatiq/aureus/loader/fileloader"
	. "github.com/dogmatiq/aureus/runner"
	"github.com/dogmatiq/aureus/test"
)

func TestRunner(t *testing.T) {
	loader := fileloader.NewLoader()

	formatJSON := func(c test.Content) ([]byte, error) {
		if len(c.Data) == 0 {
			return nil, nil
		}

		input := []byte(c.Data)

		var v any
		if err := json.Unmarshal(input, &v); err != nil {
			return nil, err
		}

		output, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return nil, err
		}

		return append(output, '\n'), nil
	}

	// expected to pass
	{
		test, err := loader.Load("testdata/pass")
		if err != nil {
			t.Fatal(err)
		}

		runner := &NativeRunner{
			Output: formatJSON,
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
			Output: formatJSON,
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
