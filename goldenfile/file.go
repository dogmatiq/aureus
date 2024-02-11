package goldenfile

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/test"
)

// File is a file that plays some role within a test.
type File struct {
	TestName  string
	IsInput   bool
	IsSkipped bool
	Content   test.Content
}

// A FileLoader loads test files from the filesystem.
//
// If the file is not part of any test, ok is false.
type FileLoader func(fsys fs.FS, name string) (f File, ok bool, err error)

// LoadFile is the default [FileLoader] implementation.
func LoadFile(fsys fs.FS, name string) (File, bool, error) {
	base := path.Base(name)
	atoms := strings.Split(base, ".")

	for i, x := range atoms {
		if i == 0 {
			// We don't look for the input/output marker in the first atom, as
			// there must be at least one atom before it from which we deduce
			// the test name.
			continue
		}

		isInput := strings.EqualFold(x, "input")
		isOutput := strings.EqualFold(x, "output")

		if !isInput && !isOutput {
			continue
		}

		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return File{}, false, err
		}

		language := ""
		if n := len(atoms) - 1; i < n {
			language = atoms[n]
		}

		testName := strings.Join(atoms[:i], ".")
		testName, skip := strings.CutPrefix(testName, "_")

		return File{
			TestName:  testName,
			IsInput:   isInput,
			IsSkipped: skip,
			Content: test.Content{
				File:     name,
				Data:     string(data),
				Language: language,
			},
		}, true, nil
	}

	return File{}, false, nil
}
