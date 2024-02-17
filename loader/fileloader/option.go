package fileloader

import (
	"io/fs"
)

// LoadOption is an option that changes the behavior of a [Loader].
type LoadOption func(*loadOptions)

type loadOptions struct {
	FS          fs.FS
	Recurse     bool
	LoadContent ContentLoader
}

// WithRecursion if a [LoadOption] that enables or disables recursive scanning
// of sub-directories.
//
// Recursion is enabled by default.
func WithRecursion(on bool) LoadOption {
	return func(opts *loadOptions) {
		opts.Recurse = on
	}
}

// WithFS is a [LoadOption] that configures an alternative filesystem to use
// when loading tests.
func WithFS(f fs.FS) LoadOption {
	return func(opts *loadOptions) {
		opts.FS = f
	}
}

// WithFileLoader is a [LoadOption] that configures an alternative [FileLoader]
// used to identify test files and load their content.
func WithFileLoader(load ContentLoader) LoadOption {
	return func(opts *loadOptions) {
		opts.LoadContent = load
	}
}
