package aureus_test

import (
	"encoding/json"
	"errors"
	"io"
	"testing"

	. "github.com/dogmatiq/aureus"
	"github.com/dogmatiq/aureus/test"
)

func prettyPrint(
	w io.Writer,
	in test.Content,
	out test.ContentMetaData,
) error {
	if in.Language != "json" || out.Language != "json" {
		return errors.New("the pretty-printer can only produce JSON output")
	}

	var v any
	if err := json.Unmarshal(
		[]byte(in.Data),
		&v,
	); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
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
