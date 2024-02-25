package streamdiff_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/dogmatiq/aureus/internal/diff"
	"github.com/dogmatiq/aureus/internal/streamdiff"
	"golang.org/x/tools/txtar"
)

func clean(text []byte) []byte {
	text = bytes.ReplaceAll(text, []byte("$\n"), []byte("\n"))
	text = bytes.TrimSuffix(text, []byte("^D\n"))
	return text
}

func Test(t *testing.T) {
	files, _ := filepath.Glob("testdata/*.txt")
	if len(files) == 0 {
		t.Fatalf("no testdata")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			a, err := txtar.ParseFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if len(a.Files) != 3 || a.Files[2].Name != "diff" {
				t.Fatalf("%s: want three files, third named \"diff\"", file)
			}

			var out bytes.Buffer
			if _, err := streamdiff.Diff(
				&out,
				a.Files[0].Name,
				bytes.NewReader(clean(a.Files[0].Data)),
				a.Files[1].Name,
				bytes.NewReader(clean(a.Files[1].Data)),
				10,
			); err != nil {
				t.Fatal(err)
			}

			got := out.Bytes()
			want := clean(a.Files[2].Data)

			if !bytes.Equal(got, want) {
				t.Logf("=== GOT ===\n%s", got)
				t.Logf("=== WANT ===\n%s", want)
				t.Fatalf("=== DIFF ===\n%s", diff.Diff("GOT", got, "WANT", want))
			}
		})
	}
}
