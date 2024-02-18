package diriter

import (
	"io/fs"
	"path"
	"strings"
)

// Each calls dir or file for each entry in the directory at the given path.
//
// It ignores any files that begin with a dot. If recurse if false, dir is never
// called.
func Each(
	fsys fs.FS,
	recurse bool,
	dirPath string,
	dir, file func(string) error,
) error {
	entries, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := path.Join(dirPath, entry.Name())

		var err error
		if !entry.IsDir() {
			err = file(entryPath)
		} else if recurse {
			err = dir(entryPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
