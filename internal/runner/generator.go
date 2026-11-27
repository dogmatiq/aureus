package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/dogmatiq/aureus/internal/test"
)

// OutputGenerator produces the output of a specific test.
type OutputGenerator[T TestingT[T]] func(T, Assertion) error

// Assertion is the content supplied to an [OutputGenerator] for one comparison.
type Assertion struct {
	Fixtures []Fixture
	Input    Input
	Output   Output
}

// Fixture is an interface for content used to prepare a test.
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

type content struct {
	io.Reader
	meta test.ContentMetaData
}

func (c *content) Name() string {
	return c.meta.Name
}

func (c *content) Language() string {
	return c.meta.Language
}

func (c *content) Attributes() map[string]string {
	return c.meta.Attributes
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

func generateOutput[T TestingT[T]](
	t T,
	gen OutputGenerator[T],
	a test.Assertion,
) (_ *os.File, err error) {
	f, err := os.CreateTemp("", "aureus-")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary file: %w", err)
	}
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	assertion := Assertion{
		Input: &content{
			Reader: bytes.NewReader(a.Input.Data),
			meta:   a.Input.ContentMetaData,
		},
		Output: &output{
			Writer: f,
			meta:   a.Output.ContentMetaData,
		},
	}

	for _, fixture := range a.Fixtures {
		assertion.Fixtures = append(
			assertion.Fixtures,
			&content{
				Reader: bytes.NewReader(fixture.Data),
				meta:   fixture.ContentMetaData,
			},
		)
	}

	if err := gen(t, assertion); err != nil {
		return nil, fmt.Errorf("unable to generate output: %w", err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("unable to seek to beginning of output file: %w", err)
	}

	return f, nil
}
