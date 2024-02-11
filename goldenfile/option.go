package goldenfile

import (
	"io/fs"
)

// LoadOption is an option that changes the behavior of a [Loader].
type LoadOption func(*loadOptions)

type loadOptions struct {
	FS       fs.FS
	Recurse  bool
	IsOutput OutputPredicate
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

// WithOutputPredicate is a [LoadOption] that configures an alternative
// predicate to use when determining whether a file is an output file.
//
// [IsOutputFile] is used by default.
func WithOutputPredicate(p OutputPredicate) LoadOption {
	return func(opts *loadOptions) {
		opts.IsOutput = p
	}
}
