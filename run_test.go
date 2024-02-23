package aureus_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/dogmatiq/aureus"
)

func prettyPrint(
	in aureus.Input,
	out aureus.Output,
) error {
	if in.Language() != "json" || out.Language() != "json" {
		return errors.New("the pretty-printer can only produce JSON output")
	}

	var v any
	dec := json.NewDecoder(in)
	if err := dec.Decode(&v); err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func TestRun_flatFile(t *testing.T) {
	aureus.Run(t, prettyPrint)
}

func TestRun_readme(t *testing.T) {
	aureus.Run(
		t,
		prettyPrint,
		aureus.WithDir("."),
		aureus.WithRecursion(false),
	)
}
