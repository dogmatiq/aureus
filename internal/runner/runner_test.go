package runner_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	"github.com/dogmatiq/aureus/internal/loader/markdownloader"
	"github.com/dogmatiq/aureus/internal/runner"
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
			GenerateOutput: func(
				_ *testing.T,
				a runner.Assertion,
			) error {
				return prettyPrint(a.Input, a.Output)
			},
			BlessStrategy: BlessDisabled,
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
				_ *testingT,
				a runner.Assertion,
			) error {
				return prettyPrint(a.Input, a.Output)
			},
			BlessStrategy: BlessDisabled,
		}

		x := &testingT{T: t}
		runner.Run(x, tst)

		for _, leaf := range x.leaves() {
			if !leaf.Failed() {
				x.Errorf("expected %q to fail", leaf.Name())
			}
		}
	}

	// expected to be blessed without changing the surrounding fixtures
	{
		dir := t.TempDir()
		if err := os.CopyFS(dir, os.DirFS("testdata/bless")); err != nil {
			t.Fatal(err)
		}

		loader := markdownloader.NewLoader()
		tst, err := loader.Load(dir)
		if err != nil {
			t.Fatal(err)
		}

		runner := &Runner[*testing.T]{
			GenerateOutput: func(
				_ *testing.T,
				a runner.Assertion,
			) error {
				return prettyPrint(a.Input, a.Output)
			},
			BlessStrategy: BlessEnabled,
		}

		runner.Run(t, tst)

		actualFile := filepath.Join(dir, "test.md")
		actual, err := os.ReadFile(actualFile)
		if err != nil {
			t.Fatal(err)
		}

		expectFile := "testdata/bless/.expect"
		expect, err := os.ReadFile(expectFile)
		if err != nil {
			t.Fatal(err)
		}

		if d := diff.Diff(expectFile, expect, actualFile, actual); d != nil {
			t.Fatal(string(d))
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

func prettyPrint(
	in runner.Input,
	out runner.Output,
) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
