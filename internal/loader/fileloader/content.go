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

		var group *loader.Group
		if i == 0 {
			// If we're in the first "atom" of the filename it means there is no
			// named portion (the filename starts with the input/output marker),
			// but we still want to group the inputs and outputs into a test
			// matrix.
			group = loader.UnnamedGroup()
		} else {
			group = loader.NamedGroup(
				strings.Join(atoms[:i], "."),
			)
		}

		return loader.Content{
			Language: lang,
			Data:     data,
			Group:    group,
			Role:     role,
		}, nil
	}

	return loader.Content{}, nil
}
