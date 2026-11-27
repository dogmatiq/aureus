package aureus_test

import (
	"encoding/json"
	"testing"

	"github.com/dogmatiq/aureus"
)

func prettyPrint(
	t *testing.T,
	a aureus.Assertion,
) error {
	if a.Input.Language() != "json" || a.Output.Language() != "json" {
		t.Fatal("the pretty-printer can only produce JSON output")
	}

	format := struct {
		Indent string `json:"indent"`
	}{
		Indent: "  ",
	}

	for _, fixture := range a.Fixtures {
		if fixture.Language() != "json" || fixture.Attributes()["role"] != "format" {
			t.Fatal("the pretty-printer only accepts JSON format fixtures")
		}

		if err := json.NewDecoder(fixture).Decode(&format); err != nil {
			return err
		}
	}

	var v any
	dec := json.NewDecoder(a.Input)
	if err := dec.Decode(&v); err != nil {
		return err
	}

	enc := json.NewEncoder(a.Output)
	enc.SetIndent("", format.Indent)
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
