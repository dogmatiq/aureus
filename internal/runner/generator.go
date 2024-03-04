package runner

import (
	"io"

	"github.com/dogmatiq/aureus/internal/test"
)

// OutputGenerator produces the output of a specific test.
type OutputGenerator[T TestingT[T]] func(T, Input, Output) error

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
	meta test.ContentMetaData
}

func (i *input) Language() string {
	return i.meta.Language
}

func (i *input) Attributes() map[string]string {
	return i.meta.Attributes
}

type output struct {
	io.Writer
	meta test.ContentMetaData
}

func (o *output) Language() string {
	return o.meta.Language
}

func (o *output) Attributes() map[string]string {
	return o.meta.Attributes
}
