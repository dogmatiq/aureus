package loadertest

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/test"
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
			var actual []byte
			test, err := load(dir)
			if err != nil {
				actual = []byte(err.Error())
			} else {
				actual = RenderTest(test)
			}

			expectFile := filepath.Join(dir, ".expect")
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
		})
	}
}
