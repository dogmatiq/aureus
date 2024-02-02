package goldenfile

import (
	"io/fs"
)

// LoadOption is an option that changes the behavior of a [Loader].
type LoadOption func(*loadOptions)

type loadOptions struct {
	FileSystem   fs.FS
	Recursive    bool
	IsGoldenFile Predicate
}

// WithRecursion if a [LoadOption] that enables or disables recursive scanning
// of sub-directories.
//
// Recursion is enabled by default.
func WithRecursion(on bool) LoadOption {
	return func(opts *loadOptions) {
		opts.Recursive = on
	}
}
