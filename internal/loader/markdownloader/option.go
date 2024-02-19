package markdownloader

import (
	"io/fs"

	"github.com/yuin/goldmark/parser"
)

// LoadOption is an option that changes the behavior of a [Loader].
type LoadOption func(*loadOptions)

type loadOptions struct {
	FS          fs.FS
	Recurse     bool
	LoadContent ContentLoader
	Parser      parser.Parser
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

// WithContentLoader is a [LoadOption] that configures an alternative [ContentLoader]
// used to load content from code blocks.
func WithContentLoader(load ContentLoader) LoadOption {
	return func(opts *loadOptions) {
		opts.LoadContent = load
	}
}

// WithParser is a [LoadOption] that configures an alternative Markdown parser
// to use when loading tests.
func WithParser(p parser.Parser) LoadOption {
	return func(opts *loadOptions) {
		opts.Parser = p
	}
}
