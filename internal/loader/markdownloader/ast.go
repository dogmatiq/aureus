package markdownloader

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// sourceOf returns the Markdown source for n.
func sourceOf(n ast.Node, source []byte) string {
	text := string(n.Text(source))
	lines := n.Lines()

	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		text += string(line.Value(source))
	}

	return text
}

// lineNumberOf returns the first line number of n.
func lineNumberOf(n ast.Node, source []byte) int {
	i := n.Lines().At(0).Start
	return bytes.Count(source[:i], []byte("\n"))
}
