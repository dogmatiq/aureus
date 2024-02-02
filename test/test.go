package test

// Test is a (possibly nested) test.
type Test struct {
	Name       string
	Flags      FlagSet
	Origin     Origin
	SubTests   []Test      `json:",omitempty"`
	Assertions []Assertion `json:",omitempty"`
}
