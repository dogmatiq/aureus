package aureus_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/dogmatiq/aureus"
)

func TestRunner(t *testing.T) {
	loader := &Loader{
		FS: os.DirFS("testdata/runner"),
	}

	formatJSON := func(a EqualAssertion) (string, error) {
		input := []byte(a.Input)

		var v any
		if err := json.Unmarshal(input, &v); err != nil {
			return "", err
		}

		output, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return "", err
		}

		return string(output) + "\n", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// expected to pass
	{
		test, err := loader.Load(ctx, "pass")
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
		test, err := loader.Load(ctx, "fail")
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
