package goldenfile

import (
	"path"
	"strings"
)

// Predicate is a predicate that determines whether a file is a
// "golden file", that is, a file that contains the expected output for some set
// of inputs.
//
// If ok is true, the file is a golden-file and p is an [InputFilePredicate]
// that that matches input files that are expected to produce output equal to
// the content of the golden file.
type Predicate func(filename string) (p InputPredicate, ok bool)

// InputPredicate is a predicate that determines whether a file contains input
// data for a test.
type InputPredicate func(filename string) bool

// DefaultPredicate is the default implementation of an [OutputPredicate].
//
// It matches any files with a ".au" extension, or a ".au.*" extension.
//
// The returned [InputPredicate] matches any files with names that begin with
// the name of the output file. It ignores any leading underscores and the
// extension of the potential input file.
func DefaultPredicate(filename string) (InputPredicate, bool) {
	prefix, ext := splitExtension(filename)
	if ext == "" {
		return nil, false
	}

	if ext != ".au" {
		prefix, ext = splitExtension(prefix)
		if ext != ".au" {
			return nil, false
		}
	}

	prefix = strings.TrimPrefix(prefix, "_")

	return func(filename string) (ok bool) {
		filename = strings.TrimPrefix(filename, "_")

		// somefile.au[.*] -> somefile
		if filename == prefix {
			return true
		}

		// somefile.au[.*] -> somefile.*
		return strings.HasPrefix(filename, prefix+".")
	}, true
}

func splitExtension(filename string) (string, string) {
	ext := path.Ext(filename)
	if ext == "" {
		return filename, ""
	}
	return strings.TrimSuffix(filename, ext), ext
}
