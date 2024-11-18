package loader

import (
	"fmt"
	"path/filepath"
)

// NoInputsError is an error that occurs when a test cannot be built because it
// has at least one input but no outputs.
type NoInputsError struct {
	Outputs []ContentEnvelope
}

func (e NoInputsError) Error() string {
	return fmt.Sprintf("output loaded from %s has no inputs", location(e.Outputs[0], true))
}

// NoOutputsError is an error that occurs when a test cannot be built because it
// has at least one output but no inputs.
type NoOutputsError struct {
	Inputs []ContentEnvelope
}

func (e NoOutputsError) Error() string {
	return fmt.Sprintf("input loaded from %s has no outputs", location(e.Inputs[0], true))
}

// location returns a string that describes the location of the given content.
func location(env ContentEnvelope, qualified bool) string {
	file := env.File
	if !qualified {
		file = filepath.Base(file)
	}

	if env.IsEntireFile() {
		return file
	}

	return fmt.Sprintf("%s:%d", file, env.Line)
}
