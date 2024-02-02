package goldenfile_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreyvit/diff"
	. "github.com/dogmatiq/aureus/goldenfile"
)

func TestLoader(t *testing.T) {
	loader := NewLoader()

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		dir := filepath.Join("testdata", e.Name())

		t.Run(e.Name(), func(t *testing.T) {
			test, err := loader.Load(dir)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := json.MarshalIndent(test, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			expect, err := os.ReadFile(filepath.Join(dir, ".expect.json"))
			if err != nil {
				t.Fatal(err)
			}

			expect = bytes.TrimSpace(expect)
			actual = bytes.TrimSpace(actual)

			if !bytes.Equal(expect, actual) {
				t.Fatalf(
					diff.LineDiff(
						string(expect),
						string(actual),
					),
				)
			}
		})
	}
}

func TestWithRecursion(t *testing.T) {
	loader := NewLoader()

	test, err := loader.Load("testdata/nested-suite", WithRecursion(false))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expect, err := os.ReadFile("testdata/nested-suite/.expect.no-recursion.json")
	if err != nil {
		t.Fatal(err)
	}

	expect = bytes.TrimSpace(expect)
	actual = bytes.TrimSpace(actual)

	if !bytes.Equal(expect, actual) {
		t.Fatalf(
			diff.LineDiff(
				string(expect),
				string(actual),
			),
		)
	}
}
