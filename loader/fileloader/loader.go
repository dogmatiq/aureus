package fileloader

import (
	"path"
	"strings"

	"github.com/dogmatiq/aureus/internal/rootfs"
	"github.com/dogmatiq/aureus/loader"
	"github.com/dogmatiq/aureus/loader/internal/diriter"
	"github.com/dogmatiq/aureus/test"
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
	return loadDir(opts, dir)
}

func loadDir(
	opts loadOptions,
	dirPath string,
) (test.Test, error) {
	var builder loader.TestBuilder

	err := diriter.Each(
		opts.FS,
		opts.Recurse,
		dirPath,
		func(dirPath string) error {
			t, err := loadDir(opts, dirPath)
			if err != nil {
				return err
			}
			builder.AddTest(t)
			return nil
		},
		func(filePath string) error {
			c, err := loadFile(opts, filePath)
			if err != nil {
				return err
			}
			builder.AddContent(c)
			return nil
		},
	)
	if err != nil {
		return test.Test{}, err
	}

	name := path.Base(dirPath)
	name, skip := strings.CutPrefix(name, "_")

	subTests, err := builder.Build()
	if err != nil {
		return test.Test{}, err
	}

	return test.New(
		name,
		test.WithSkip(skip),
		test.WithSubTests(subTests...),
	), nil
}

func loadFile(opts loadOptions, filePath string) (loader.ContentEnvelope, error) {
	f, err := opts.FS.Open(filePath)
	if err != nil {
		return loader.ContentEnvelope{}, err
	}
	defer f.Close()

	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")

	c, err := opts.LoadContent(name, f)
	if err != nil {
		return loader.ContentEnvelope{}, err
	}

	return loader.ContentEnvelope{
		File:    filePath,
		Skip:    skip,
		Content: c,
	}, nil
}
