package markdownloader

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// linesOf returns the text within n.
func linesOf(n ast.Node, source []byte) string {
	return string(n.Lines().Value(source))
}

var newline = []byte("\n")

// locationOf returns the location of n within source.
func locationOf(n ast.Node, source []byte) (line, begin, end int) {
	lines := n.Lines()
	count := lines.Len()

	begin = lines.At(0).Start
	end = lines.At(count - 1).Stop
	line = bytes.Count(source[:begin], newline)

	return line, begin, end
}
