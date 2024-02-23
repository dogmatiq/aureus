package fileloader

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/internal/loader"
	"github.com/dogmatiq/aureus/internal/rootfs"
	"github.com/dogmatiq/aureus/internal/test"
)

// Loader loads [test.Test] values from files containing test inputs and
// expected outputs.
type Loader struct {
	options loadOptions
}

// NewLoader returns a new [Loader], which loads golden file tests from the
// filesystem.
func NewLoader(options ...LoadOption) *Loader {
	l := &Loader{
		options: loadOptions{
			FS:          rootfs.FS,
			Recurse:     true,
			LoadContent: LoadContent,
		},
	}

	for _, opt := range options {
		opt(&l.options)
	}

	return l
}

// Load returns a test built from files in the given directory.
//
// Any directory or file that begins with an underscore produces a test that is
// marked as skipped.
func (l *Loader) Load(dir string, options ...LoadOption) (test.Test, error) {
	opts := l.options
	for _, opt := range options {
		opt(&opts)
	}

	return loader.LoadDir(
		opts.FS,
		dir,
		opts.Recurse,
		func(builder *loader.TestBuilder, fsys fs.FS, filePath string) error {
			f, err := fsys.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			name := path.Base(filePath)
			name, skip := strings.CutPrefix(name, "_")

			c, err := opts.LoadContent(name, f)
			if err != nil {
				return err
			}

			return builder.AddContent(
				loader.ContentEnvelope{
					File:    filePath,
					Skip:    skip,
					Content: c,
				},
			)
		},
	)
}
