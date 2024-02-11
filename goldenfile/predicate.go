package goldenfile

import (
	"slices"
	"strings"
)

// An OutputPredicate is a function that determines whether a file is to be
// treated as an "output file", that is, a file that contains the expected
// output for some set of inputs.
//
// If filename refers to an output file, ok is true and p is an [InputPredicate]
// that matches the "input files" that are expected to produce output equal to
// the content of the output file.
type OutputPredicate func(filename string) (p InputPredicate, ok bool)

// An InputPredicate is a function that determines whether a file is to be
// treated as an "input file" for a specific "output file".
type InputPredicate func(filename string) bool

// IsOutputFile is the default [OutputPredicate] implementation.
//
// It matches any filenames with an ".output" part, that is, either the
// extension ".output", or with ".output." appearing in the filename.
//
// The returned [InputPredicate] matches filenames with the same prefix as the
// output file, up to the first occurrence of ".output".
func IsOutputFile(filename string) (InputPredicate, bool) {
	if want, ok := hasAtom(filename, "output"); ok {
		return func(filename string) bool {
			got, ok := hasAtom(filename, "input")
			return ok && slices.Equal(got, want)
		}, true
	}
	return nil, false
}

// hasAtom returns the atoms (dot-separated components) of filename that precede
// the first occurrence of a.
//
// It returns false if a is the first atom, or does not appear in the filename
// at all.
//
// Leading underscores are ignored.
func hasAtom(filename, a string) ([]string, bool) {
	filename = strings.TrimPrefix(filename, "_")
	atoms := strings.Split(filename, ".")

	for i, x := range atoms[1:] {
		if x == a {
			return atoms[:i+1], true
		}
	}

	return nil, false
}
