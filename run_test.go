package aureus_test

import (
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/dogmatiq/aureus"
)

func prettyPrint(
	w io.Writer,
	in aureus.Content,
	out aureus.ContentMetaData,
) error {
	if in.Language != "json" || out.Language != "json" {
		return errors.New("the pretty-printer can only produce JSON output")
	}

	var v any
	if err := json.Unmarshal(in.Data, &v); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
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
