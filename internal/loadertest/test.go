package loadertest

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/test"
)

// Run executes a set of basic golden-file tests for a loader using directories
// within the testdata directory.
//
// It does not use Aureus to execute the tests, so that it can be used to test
// Aureus loader impelmentations.
func Run[O any](
	t *testing.T,
	load func(string, ...O) (test.Test, error),
) {
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

			test, err := load(dir)
			if err != nil {
				enc.Encode(struct{ Error string }{err.Error()})
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
