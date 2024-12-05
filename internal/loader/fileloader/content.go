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

	content := loader.Content{
		// If there is no prefix on the filename before the .input or .output
		// atom marker, we still want to group the inputs and outputs into a
		// test matrix.
		Group: loader.UnnamedGroup(),
	}

	for idx, atom := range atoms {
		if strings.EqualFold(atom, "input") {
			content.Role = loader.Input
		} else if strings.EqualFold(atom, "output") {
			content.Role = loader.Output
		} else {
			continue
		}

		if idx > 0 {
			group := strings.Join(atoms[:idx], ".")
			content.Group = loader.NamedGroup(group)
		}

		atoms = atoms[idx+1:]
		break
	}

	if content.Role == loader.NoRole {
		return loader.Content{}, nil
	}

	for len(atoms) != 0 {
		atom := atoms[0]
		attr, ok := strings.CutPrefix(atom, "@")
		if !ok {
			break
		}
		atoms = atoms[1:]

		if content.Attributes == nil {
			content.Attributes = make(map[string]string)
		}

		if pos := strings.Index(attr, "="); pos != -1 {
			content.Attributes[attr[:pos]] = attr[pos+1:]
		} else {
			content.Attributes[attr] = ""
		}
	}

	content.Language = strings.Join(atoms, ".")

	var err error
	content.Data, err = io.ReadAll(f)
	if err != nil {
		return loader.Content{}, err
	}

	return content, nil
}
