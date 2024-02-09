package goldenfile

import (
	"slices"
	"strings"
)

// A Predicate is a function that determines whether a file is an "output file",
// that is, a file that contains the expected output for some set of inputs.
//
// If ok is true, the file is an output file and p is an [InputFilePredicate]
// that that matches input files that are expected to produce the same output as
// the contents of the output file.
type Predicate func(filename string) (p InputPredicate, ok bool)

// InputPredicate is a predicate that determines whether a file contains input
// data for a test.
type InputPredicate func(filename string) bool

// DefaultPredicate is the default [Predicate] implementation.
//
// It matches any filenames with an ".output" part, that is, either the
// extension ".output", or with ".output." appearing in the filename.
//
// The returned [InputPredicate] matches filenames with the same prefix as the
// output file, up to the first occurrence of ".output".
func DefaultPredicate(filename string) (InputPredicate, bool) {
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
