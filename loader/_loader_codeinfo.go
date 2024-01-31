package aureus

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/yuin/goldmark/ast"
	"golang.org/x/net/html"
)

const (
	defaultLanguage    = "text"
	attributeNamespace = "au"
)

// codeInfo is the result of parsing the "info line" of a fenced code block.
type codeInfo struct {
	Language    string
	IsAssertion bool
	Skip        bool
}

// parseCodeInfo parses the "info line" (the line indicating the block's
// language) of a code block.
func (l *Loader) parseCodeInfo(n *ast.FencedCodeBlock, source []byte) (codeInfo, error) {
	info := codeInfo{
		Language: defaultLanguage,
	}

	if n.Info == nil {
		return info, nil
	}

	r := io.MultiReader(
		strings.NewReader(`<html `),
		bytes.NewReader(n.Info.Segment.Value(source)),
		strings.NewReader(">"),
	)

	node, err := html.Parse(r)
	if err != nil {
		return codeInfo{}, err
	}

	for i, attr := range node.FirstChild.Attr {
		key, ok := strings.CutPrefix(attr.Key, attributeNamespace+":")
		if !ok {
			if i == 0 && attr.Val == "" {
				info.Language = attr.Key
			}
			continue
		}

		switch key {
		case "assertion":
			info.IsAssertion = true
			if attr.Val != "" {
				return codeInfo{}, fmt.Errorf("unexpected value for %q attribute", attr.Key)
			}
		case "skip":
			switch attr.Val {
			case "", "true":
				info.Skip = true
			case "false":
				info.Skip = false
			default:
				return codeInfo{}, fmt.Errorf("unexpected value for %q attribute; %q", attr.Key, attr.Val)
			}
		}
	}

	return info, nil
}
