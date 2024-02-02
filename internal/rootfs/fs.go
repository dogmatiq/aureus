package rootfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

var underlying = os.DirFS("/")

// FS is the root OS file system.
//
// It respects the current working directory when referencing relative paths.
var FS fs.FS = rootFS{}

type rootFS struct{}

func (rootFS) Open(name string) (fs.File, error) {
	name, err := normalizePath(name)
	if err != nil {
		return nil, err
	}
	return underlying.Open(name)
}

func (rootFS) ReadFile(name string) ([]byte, error) {
	name, err := normalizePath(name)
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(underlying, name)
}

func (rootFS) Stat(name string) (fs.FileInfo, error) {
	name, err := normalizePath(name)
	if err != nil {
		return nil, nil
	}
	return fs.Stat(underlying, name)
}

func normalizePath(name string) (string, error) {
	name, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	return filepath.Rel("/", name)
}
