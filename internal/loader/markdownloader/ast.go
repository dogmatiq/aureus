package markdownloader

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// linesOf returns the text within n.
func linesOf(n ast.Node, source []byte) string {
	return string(n.Lines().Value(source))
}

// lineNumberOf returns the first line number of n.
func lineNumberOf(n ast.Node, source []byte) int {
	i := n.Lines().At(0).Start
	return bytes.Count(source[:i], []byte("\n"))
}
