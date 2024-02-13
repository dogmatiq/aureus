package goldenfile_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/dogmatiq/aureus/goldenfile"
	"github.com/dogmatiq/aureus/internal/diff"
)

func TestLoader(t *testing.T) {
	loader := NewLoader()

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}

		dir := filepath.Join("testdata", e.Name())

		t.Run(e.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "  ")

			test, err := loader.Load(dir)
			if err != nil {
				enc.Encode(
					struct {
						Error string
					}{
						err.Error(),
					},
				)
			} else if err := enc.Encode(test); err != nil {
				t.Fatal(err)
			}

			expectFile := filepath.Join(dir, ".expect.json")
			expect, err := os.ReadFile(expectFile)
			if err != nil {
				t.Fatal(err)
			}

			expect = bytes.TrimSpace(expect)
			actual := bytes.TrimSpace(buf.Bytes())

			if d := diff.Diff(
				expectFile, expect,
				"actual.json", actual,
			); d != nil {
				t.Fatal(string(d))
			}
		})
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
