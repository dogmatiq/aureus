package loadertest

import (
	"bytes"
	"fmt"

	"github.com/dogmatiq/aureus/internal/test"
)

// RenderTest returns a byte slice containing a human-readable representation of
// a test.
func RenderTest(t test.Test) []byte {
	var w bytes.Buffer
	w.WriteString("test")

	if t.Name != "" {
		fmt.Fprintf(&w, " %q", t.Name)
	}

	if t.Skip {
		w.WriteString(" [skipped]")
	}

	w.WriteString(" {\n")

	for _, s := range t.SubTests {
		indent(&w, RenderTest(s))
	}

	for _, a := range t.Assertions {
		indent(&w, renderAssertion(a))
	}

	w.WriteString("}")

	return w.Bytes()
}

func renderAssertion(a test.Assertion) []byte {
	var w bytes.Buffer
	w.WriteString("assertion {\n")
	indent(&w, renderContent("input", a.Input))
	indent(&w, renderContent("output", a.Output))
	w.WriteString("}")
	return w.Bytes()
}

func renderContent(label string, c test.Content) []byte {
	var w bytes.Buffer

	w.WriteString(label)
	w.WriteByte(' ')

	if c.IsEntireFile() {
		fmt.Fprintf(&w, "%q", c.File)
	} else {
		fmt.Fprintf(&w, `"%s:%d"`, c.File, c.Line)
	}

	w.WriteString(" {\n")

	if c.Language != "" {
		fmt.Fprintf(&w, "    lang = %q\n", c.Language)
	}
	fmt.Fprintf(&w, "    data = %q\n", string(c.Data))

	w.WriteString("}")

	return w.Bytes()
}

func indent(w *bytes.Buffer, data []byte) {
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		w.WriteString("    ")
		w.Write(line)
		w.WriteString("\n")
	}
}
