package test

import (
	"fmt"
	"path"
	"strings"
)

// Test is a (possibly nested) test.
type Test struct {
	Name      string
	Flags     FlagSet   `json:",omitempty"`
	SubTests  []Test    `json:",omitempty"`
	Assertion Assertion `json:",omitempty"`
}

// New creates a new [Test].
//
// It returns the test and the set of flags that should be inherited by any
// sub-tests.
func New(options ...Option) (Test, FlagSet) {
	var opts testOptions
	for _, opt := range options {
		opt(&opts)
	}

	opts.Flags.Add(opts.InheritedFlags)

	if opts.Flags.Has(FlagSkipped) {
		opts.InheritedFlags.Add(FlagAncestorSkipped)
	}

	return opts.Test, opts.InheritedFlags
}

// Option is an option that controls how a test is created by [New].
type Option func(*testOptions)

type testOptions struct {
	Test
	InheritedFlags FlagSet
}

// WithNameFromPath is a [TestOption] that sets the name and flags of a
// test based on its source path, which may be either a file or a directory.
func WithNameFromPath(p string) Option {
	name := path.Base(p)
	name, skip := strings.CutPrefix(name, "_")

	return func(opts *testOptions) {
		opts.Name = name
		if skip {
			opts.Flags.Add(FlagSkipped)
		}
	}
}

// WithName is a [TestOption] that sets the name of a test.
func WithName(format string, args ...any) Option {
	return func(opts *testOptions) {
		opts.Name = fmt.Sprintf(format, args...)
	}
}

// WithFlag is a [TestOption] that sets the given flags on a test.
func WithFlag(f FlagLike) Option {
	return func(opts *testOptions) {
		opts.Flags.Add(f)
	}
}

// WithInheritedFlags is a [TestOption] that adds the given flags to the set of
// flags that are inherited from the test's ancestors.
func WithInheritedFlags(s FlagSet) Option {
	return func(opts *testOptions) {
		opts.InheritedFlags.Add(s)
	}
}

// WithAssertion is a [TestOption] that sets the assertion on the test.
func WithAssertion(a Assertion) Option {
	return func(opts *testOptions) {
		opts.Assertion = a
	}
}
