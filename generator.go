package aureus

import (
	"io"
)

// OutputGenerator produces the output of a specific test.
type OutputGenerator[T TestingT[T]] func(T, Input, Output)

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

type input struct {
	io.Reader
	lang  string
	attrs map[string]string
}

func (i *input) Language() string {
	return i.lang
}

func (i *input) Attributes() map[string]string {
	return i.attrs
}

type output struct {
	io.Writer
	lang  string
	attrs map[string]string
}

func (o *output) Language() string {
	return o.lang
}

func (o *output) Attributes() map[string]string {
	return o.attrs
}
