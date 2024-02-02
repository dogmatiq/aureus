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

// flags is a slice of all possible flags.
var flags = []Flag{
	FlagSkipped,
	FlagAncestorSkipped,
}

// FlagSet is a set of flags.
type FlagSet int

// EmptyFlagSet is a [FlagSet] containing no flags.
const EmptyFlagSet FlagSet = 0

// Has returns true if f contains the given flag.
func (s FlagSet) Has(f Flag) bool {
	return s&FlagSet(f) != 0
}

// Add adds the given flag to f.
func (s *FlagSet) Add(f Flag) {
	*s |= FlagSet(f)
}

// Remove removes the given flag from f.
func (s *FlagSet) Remove(f Flag) {
	*s &^= FlagSet(f)
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
