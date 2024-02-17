package goldenfile_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	. "github.com/dogmatiq/aureus/goldenfile"
	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/loadertest"
)

func TestLoader(t *testing.T) {
	loader := NewLoader()
	loadertest.Run(t, loader.Load)
}

func TestWithRecursion(t *testing.T) {
	loader := NewLoader()

	test, err := loader.Load("testdata/nested-directory", WithRecursion(false))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expectFile := "testdata/nested-directory/.expect.no-recursion.json"
	expect, err := os.ReadFile(expectFile)
	if err != nil {
		t.Fatal(err)
	}

	expect = bytes.TrimSpace(expect)
	actual = bytes.TrimSpace(actual)

	if d := diff.Diff(
		expectFile, expect,
		"actual.json", actual,
	); d != nil {
		t.Fatal(string(d))
	}
}
