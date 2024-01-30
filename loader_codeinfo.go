package aureus

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

const (
	defaultLanguage = "text"
	attributeKey    = "au"
)

// codeInfo is the result of parsing the "info line" of a fenced code block.
type codeInfo struct {
	Language    string
	IsAssertion bool
	Skip        bool
}

// parseCodeInfo parses the "info line" (the line indicating the block's
// language) of a code block.
func (l *Loader) parseCodeInfo(n *ast.FencedCodeBlock, source []byte) codeInfo {
	info := codeInfo{
		Language: l.DefaultLanguage,
	}
	if info.Language == "" {
		info.Language = defaultLanguage
	}

	if n.Info == nil {
		return info
	}

	text := string(n.Info.Segment.Value(source))
	fields := strings.Fields(text)

	for i, field := range fields {
		if field == attributeKey {
			info.IsAssertion = true
		} else if i == 0 {
			info.Language = field
		} else if field == "skip" && info.IsAssertion {
			info.Skip = true
		}
	}

	return info
}
