package cliflags

import (
	"flag"
)

// Flags is a struct that holds all Aureus command-line flags.
type Flags struct {
	Bless bool
}

// Get returns the Aureus command-line flags.
func Get() Flags {
	return flags
}

var flags Flags

func init() {
	flag.BoolVar(
		&flags.Bless,
		"aureus.bless",
		false,
		"replace (on disk) each failing assertion's expected output with its current output",
	)
}
