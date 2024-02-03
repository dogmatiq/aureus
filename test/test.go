package test

import (
	"fmt"
	"path"
	"strings"
)

// Test is a (possibly nested) test.
type Test struct {
	Name      string
	Flags     FlagSet
	Origin    Origin
	SubTests  []Test    `json:",omitempty"`
	Assertion Assertion `json:",omitempty"`
}

// New creates a new [Test].
//
// The name of the test is derived from the path in the origin. If the basename
// of the origin starts with an underscore, the test is marked as skipped.
//
// It returns the test and the set of flags that should be inherited by any
// sub-tests.
func New(origin Origin, options ...Option) (Test, FlagSet) {
	name := path.Base(origin.Path())
	name, skip := strings.CutPrefix(name, "_")

	t := Test{
		Name:   name,
		Origin: origin,
	}

	for _, opt := range options {
		opt(&t)
	}

	inherited := t.Flags

	if skip {
		t.Flags.Add(FlagSkipped)
		inherited.Add(FlagAncestorSkipped)
	}

	return t, inherited
}

// Option is an option that controls how a test is created by [New].
type Option func(*Test)

// WithName is a [TestOption] that overrides the default name of a [Test].
func WithName(format string, args ...any) Option {
	return func(t *Test) {
		t.Name = fmt.Sprintf(format, args...)
	}
}

// WithFlag is a [TestOption] that sets the given flags on a test.
func WithFlag(flags ...Flag) Option {
	return func(t *Test) {
		t.Flags.Add(flags...)
	}
}

// WithInheritedFlags is a [TestOption] that sets the given set of inherited
// flags on a test. Inherited flags are those inherited by the test's ancestors.
func WithInheritedFlags(inherited FlagSet) Option {
	return func(t *Test) {
		t.Flags = inherited
	}
}

// WithSubTest is a [TestOption] that adds sub-tests to the test.
func WithSubTest(subs ...Test) Option {
	return func(t *Test) {
		t.SubTests = append(t.SubTests, subs...)
	}
}

// WithAssertion is a [TestOption] that sets the assertion on the test.
func WithAssertion(a Assertion) Option {
	return func(t *Test) {
		t.Assertion = a
	}
}
