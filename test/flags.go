package test

import "strings"

// Flag is a bitmask of flags describing a [Node].
type Flag int

const (
	// FlagSkipped indicates that the test should be skipped.
	FlagSkipped Flag = 1 << iota

	// FlagAncestorSkipped indicates that the test should be skipped because one
	// of its ancestors was skipped.
	FlagAncestorSkipped
)

func (f Flag) String() string {
	switch f {
	case FlagSkipped:
		return "skipped"
	case FlagAncestorSkipped:
		return "ancestor-skipped"
	default:
		return "unknown"
	}
}

func (f Flag) flags() FlagSet {
	return FlagSet(f)
}

// flags is a slice of all possible flags.
var flags = []Flag{
	FlagSkipped,
	FlagAncestorSkipped,
}

// FlagSet is a set of flags.
type FlagSet int

// EmptyFlagSet is a [FlagSet] containing no flags.
const EmptyFlagSet FlagSet = 0

// Has returns true if s contains the given flag.
func (s FlagSet) Has(f FlagLike) bool {
	return s&f.flags() != 0
}

// Add adds the given flags to s.
func (s *FlagSet) Add(f FlagLike) {
	*s |= f.flags()
}

// Remove removes the given flags from s.
func (s *FlagSet) Remove(f FlagLike) {
	*s &^= f.flags()
}

func (s FlagSet) String() string {
	if s == EmptyFlagSet {
		return "-"
	}

	var names []string
	for _, f := range flags {
		if s.Has(f) {
			names = append(names, f.String())
		}
	}
	return strings.Join(names, " | ")
}

// DapperString returns the string used to represent s in
// [github.com/dogmatiq/dapper] output.
func (s FlagSet) DapperString() string {
	return s.String()
}

func (s FlagSet) flags() FlagSet {
	return s
}

// FlagLike is a type that can be converted to a [FlagSet].
type FlagLike interface {
	flags() FlagSet
}
