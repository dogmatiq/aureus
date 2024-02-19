package aureus_test

import (
	"encoding/json"
	"io"
	"testing"

	. "github.com/dogmatiq/aureus"
	"github.com/dogmatiq/aureus/test"
)

func prettyPrint(input test.Content, output io.Writer) error {
	if len(input.Data) == 0 {
		return nil
	}

	var v any
	if err := json.Unmarshal(
		[]byte(input.Data),
		&v,
	); err != nil {
		return err
	}

	enc := json.NewEncoder(output)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func TestRun_flatFile(t *testing.T) {
	Run(t, prettyPrint)
}

func TestRun_readme(t *testing.T) {
	Run(
		t,
		prettyPrint,
		WithDir("."),
		WithRecursion(false),
	)
}
