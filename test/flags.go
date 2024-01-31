package test

// Flags is a bitmask of flags describing a [Node].
type Flags int

const (
	// FlagNone is the default value for [Flags].
	FlagNone Flags = 0

	// FlagSkipped indicates that the test should be skipped.
	FlagSkipped Flags = 1 << iota

	// FlagAncestorSkipped indicates that the test should be skipped because one
	// of its ancestors was skipped.
	FlagAncestorSkipped
)

// Has returns true if f contains the given flag.
func (f *Flags) Has(flag Flags) bool {
	return *f&flag != 0
}

// Add adds the given flag to f.
func (f *Flags) Add(flag Flags) {
	*f |= flag
}

// Remove removes the given flag from f.
func (f *Flags) Remove(flag Flags) {
	*f &^= flag
}
