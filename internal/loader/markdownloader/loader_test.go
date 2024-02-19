package markdownloader_test

import (
	"testing"

	"github.com/dogmatiq/aureus/internal/loader/internal/loadertest"
	. "github.com/dogmatiq/aureus/internal/loader/markdownloader"
)

func TestLoader(t *testing.T) {
	loader := NewLoader()
	loadertest.Run(t, loader.Load)
}
