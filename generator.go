package aureus

import (
	"io"
)

// OutputGenerator produces the output of a specific test.
type OutputGenerator[T TestingT[T]] func(T, Assertion) error

// Assertion is the content supplied to an [OutputGenerator] for one comparison.
type Assertion struct {
	// Fixtures contains the content used to prepare the assertion, in document
	// or filename order. Each assertion receives fresh readers, including when
	// a group forms a test matrix.
	Fixtures []Fixture

	// Input is the value from which the generator produces output.
	Input Input

	// Output receives the actual output to compare with the expected output.
	Output Output
}

// Fixture is content used to prepare an assertion.
//
// Aureus supplies fixtures to the generator; it does not interpret or execute
// them. The generator can use their languages and attributes to prepare state.
type Fixture interface {
	io.Reader

	// Language returns the language of the fixture, if known.
	Language() string

	// Attributes returns loader-specific information about the fixture.
	Attributes() map[string]string
}

// Input is an interface for the input to a test.
type Input interface {
	io.Reader

	// Language returns the language of the input value, if known, e.g. "json",
	// "yaml", etc.
	Language() string

	// Attributes returns a set of key-value pairs that provide additional
	// loader-specific information about the input.
	Attributes() map[string]string
}

// Output is an interface for producing the output for a test.
type Output interface {
	io.Writer

	// Language returns the expected language of the output value, if known,
	// e.g. "json", "yaml", etc.
	Language() string

	// Attributes returns a set of key-value pairs that provide additional
	// loader-specific information about the expected output.
	Attributes() map[string]string
}
