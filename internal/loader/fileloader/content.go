package fileloader

import (
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/internal/loader"
)

// A ContentLoader is a function that returns content read from a file.
//
// name is the "effective" or "sanitized" name of the file, after any special
// characters have been removed by the loader. For example, a leading
// underscore. The f.Stat() method can be used to get the actual file name.
//
// If the returned content's role is [loader.NoRole], it is ignored.
type ContentLoader func(name string, f fs.File) (loader.Content, error)

// LoadContent is the default [ContentLoader] implementation.
func LoadContent(name string, f fs.File) (loader.Content, error) {
	base := path.Base(name)
	atoms := strings.Split(base, ".")

	for i, x := range atoms {
		if i == 0 {
			// We don't look for the input/output marker in the first atom, as
			// there must be at least one atom before it from which we deduce
			// the test name.
			continue
		}

		var role loader.ContentRole
		if strings.EqualFold(x, "input") {
			role = loader.Input
		} else if strings.EqualFold(x, "output") {
			role = loader.Output
		} else {
			continue
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return loader.Content{}, err
		}

		lang := ""
		if n := len(atoms) - 1; i < n {
			lang = atoms[n]
		}

		return loader.Content{
			Language: lang,
			Data:     string(data),
			Group:    strings.Join(atoms[:i], "."),
			Role:     role,
		}, nil
	}

	return loader.Content{}, nil
}
