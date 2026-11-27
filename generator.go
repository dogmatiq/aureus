package aureus

import (
	"io"
)

// OutputGenerator produces the output of a specific test.
type OutputGenerator[T TestingT[T]] func(T, Assertion) error

// Assertion is the content supplied to an [OutputGenerator] for one comparison.
type Assertion struct {
	// Fixtures contains the group's fixtures and any required fixtures in
	// dependency order. Each fixture appears once, with a fresh reader for each
	// assertion.
	Fixtures []Fixture

	// Input is the value from which the generator produces output.
	Input Input

	// Output receives the actual output to compare with the expected output.
	Output Output
}

// Fixture is an interface for content used to prepare a test.
//
// Aureus supplies fixtures to the generator; it does not interpret or execute
// them. The generator can use their languages and attributes to prepare state.
//
// Fixture names are local to one Markdown document or directory of flat-file
// tests. A fixture annotated with au:group (or @group in its filename) is
// supplied to every assertion in that group. Inputs, outputs and fixtures can
// also require ungrouped fixtures using au:requires (or @requires in filenames).
// Names are comma-separated.
//
// Dependencies are supplied before their dependents, following the listed
// requirements from the group's fixtures, then the input and output. Group
// fixtures are visited in declaration order. Missing required fixtures,
// references to grouped fixtures, dependency cycles and duplicate declarations
// are errors. Every fixture must be used by at least one assertion, directly or
// through another fixture. Skipping a fixture skips assertions that use it.
type Fixture interface {
	io.Reader

	// Name returns the name used to declare and reference the fixture.
	Name() string

	// Language returns the language of the fixture, if known, e.g. "json",
	// "yaml", etc.
	Language() string

	// Attributes returns a set of key-value pairs that provide additional
	// loader-specific information about the fixture.
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
