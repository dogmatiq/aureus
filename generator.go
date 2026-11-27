package aureus

import (
	"io"
)

// OutputGenerator produces the output of a specific test.
//
// Implementations prepare and clean up any state needed for an assertion.
type OutputGenerator[T TestingT[T]] func(T, Assertion) error

// Assertion is the content supplied to an [OutputGenerator] for one comparison.
type Assertion struct {
	// Fixtures contains the explicitly required fixtures in dependency order.
	// Each fixture appears once, with a fresh reader for each assertion.
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
//
// Fixture names are local to one Markdown document or directory of flat-file
// tests. Inputs, outputs and fixtures declare their dependencies using
// au:requires (or @requires in filenames). Names are comma-separated.
//
// Dependencies are supplied before their dependants, following the listed
// requirements from the input, then the output. Missing required fixtures,
// dependency cycles and duplicate declarations are errors. Unused fixtures
// are not included. Skipping a fixture skips assertions that require it,
// including through another fixture.
type Fixture interface {
	io.Reader

	// Name returns the name used to declare and reference the fixture.
	Name() string

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
