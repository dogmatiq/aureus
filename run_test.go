package aureus_test

import (
	"encoding/json"
	"testing"

	"github.com/dogmatiq/aureus"
)

func prettyPrint(
	t *testing.T,
	in aureus.Input,
	out aureus.Output,
) error {
	if in.Language() != "json" || out.Language() != "json" {
		t.Fatal("the pretty-printer can only produce JSON output")
	}

	var v any
	dec := json.NewDecoder(in)
	if err := dec.Decode(&v); err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return err
	}

	return nil
}

func TestRun_flatFile(t *testing.T) {
	aureus.Run(t, prettyPrint)
}

func TestRun_readme(t *testing.T) {
	aureus.Run(
		t,
		prettyPrint,
		aureus.FromDir("."),
		aureus.Recursive(false),
	)
}
