package fsiter

import (
	"io/fs"
	"path"
	"strings"
)

func Each[T any](
	fsys fs.FS,
	recurse bool,
	dirPath string,
	dir func(string) (T, error),
	file func(string) error,
) ([]T, error) {
	entries, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return nil, err
	}

	var results []T

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := path.Join(dirPath, entry.Name())

		if !entry.IsDir() {
			if err := file(entryPath); err != nil {
				return nil, err
			}
		} else if recurse {
			v, err := dir(entryPath)
			if err != nil {
				return nil, err
			}
			results = append(results, v)
		}
	}

	return results, nil
}
