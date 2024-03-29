package fileloader_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	. "github.com/dogmatiq/aureus/internal/loader/fileloader"
	"github.com/dogmatiq/aureus/internal/loader/internal/loadertest"
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

	expectFile := "testdata/nested-directory/.expect.no-recursion"
	expect, err := os.ReadFile(expectFile)
	if err != nil {
		t.Fatal(err)
	}

	expect = bytes.TrimSpace(expect)
	actual := bytes.TrimSpace(loadertest.RenderTest(test))

	if d := diff.Diff(
		expectFile, expect,
		"actual", actual,
	); d != nil {
		t.Fatal(string(d))
	}
}
