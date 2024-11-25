package runner

import (
	"fmt"
	"regexp"
	"strings"
)

func (r *Runner[T]) goTestCommand(t LoggerT) string {
	pkg := "./..."
	if r.PackagePath != "" {
		pkg = shellQuote(r.PackagePath)
	}

	return fmt.Sprintf(
		"go test %s -run %s -v -count 1",
		pkg,
		shellQuote(testNamePattern(t)),
	)
}

func testNamePattern(t LoggerT) string {
	atoms := strings.Split(t.Name(), "/")
	for i, atom := range atoms {
		atoms[i] = "^" + regexp.QuoteMeta(atom) + "$"
	}
	return strings.Join(atoms, "/")
}

func shellQuote(s string) string {
	var w strings.Builder
	w.WriteByte('\'')

	for _, r := range s {
		if r == '\'' {
			w.WriteString("'\"'\"'")
		} else {
			w.WriteRune(r)
		}
	}

	w.WriteByte('\'')
	return w.String()
}
