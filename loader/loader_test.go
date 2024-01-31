package loader_test

import (
	"testing"

	. "github.com/dogmatiq/aureus/loader"
)

func TestLoader(t *testing.T) {
	loader := &Loader{}

	t.Run("nested", func(t *testing.T) {
		_, err := loader.Load("testdata/nested.au.md")
		if err != nil {
			t.Fatal(err)
		}
	})
}
