package test

import (
	"path"
	"strings"
)

// Test is a (possibly nested) test.
type Test struct {
	Name       string
	Flags      FlagSet
	Origin     Origin
	SubTests   []Test      `json:",omitempty"`
	Assertions []Assertion `json:",omitempty"`
}

// New creates a new [Test] from an [Origin] and an inherited set of flags.
func New(o Origin, inherited FlagSet) (Test, FlagSet) {
	name := path.Base(o.Path())
	name, skip := strings.CutPrefix(name, "_")

	t := Test{
		Name:   name,
		Flags:  inherited,
		Origin: o,
	}

	if skip {
		t.Flags.Add(FlagSkipped)
		inherited.Add(FlagAncestorSkipped)
	}

	return t, inherited
}
