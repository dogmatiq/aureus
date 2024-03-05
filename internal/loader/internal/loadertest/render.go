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

	fmt.Fprintf(&w, "test %q", t.Name)
	if t.Skip {
		w.WriteString(" [skipped]")
	}
	w.WriteString(" {\n")

	for _, s := range t.SubTests {
		indent(&w, RenderTest(s))
	}

	var r assertionRenderer
	for _, a := range t.Assertions {
		r.w.Reset()
		a.AcceptVisitor(&r)
		indent(&w, r.w.Bytes())
	}

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

type assertionRenderer struct {
	w bytes.Buffer
}

func (r *assertionRenderer) VisitEqualAssertion(a test.EqualAssertion) {
	r.w.WriteString("equal {\n")
	indent(&r.w, renderContent("input", a.Input))
	indent(&r.w, renderContent("output", a.Output))
	r.w.WriteString("}")
}

func renderContent(label string, c test.Content) []byte {
	var w bytes.Buffer

	w.WriteString(label)
	w.WriteByte(' ')

	if c.Line == 0 {
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
