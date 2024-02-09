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
				t.Log(
					string(actual),
				)
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

func TestLoader_outputFileWithNoInputs(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load(
		"testdata/without-file-extension",
		WithPredicate(
			func(filename string) (InputPredicate, bool) {
				if _, ok := DefaultPredicate(filename); ok {
					return func(filename string) bool {
						return false
					}, true
				}
				return nil, false
			},
		),
	)
	expect := `output file "testdata/without-file-extension/test.output" has no associated input files`
	if err == nil {
		t.Fatalf("expected an error: got nil, want %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("unexpected error: got %q, want %q", err.Error(), expect)
	}
}

func TestLoader_inputFileWithNoOutputs(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load(
		"testdata/without-file-extension",
		WithPredicate(
			func(filename string) (InputPredicate, bool) {
				return nil, false
			},
		),
	)

	expect := `input file "testdata/without-file-extension/test.input" has no associated output files`
	if err == nil {
		t.Fatalf("expected an error: got nil, want %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("unexpected error: got %q, want %q", err.Error(), expect)
	}
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

	expect, err := os.ReadFile("testdata/nested-directory/.expect.no-recursion.json")
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
