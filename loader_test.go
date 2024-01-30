package aureus_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/dogmatiq/aureus"
)

func TestLoader(t *testing.T) {
	loader := &Loader{
		FS: os.DirFS("testdata/loader"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("missing input", func(t *testing.T) {
		_, err := loader.Load(ctx, "missing-input.au.md")
		if err == nil {
			t.Fatal("expected an error")
		}

		want := `assertion in missing-input.au.md on line 1: preceding code block containing the input value was not found`
		got := err.Error()
		if got != want {
			t.Errorf("expected error: got %q, want %q", got, want)
		}
	})
}
