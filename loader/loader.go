package loader

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/test"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// FileExtension is the file extension used to identify Aureus Markdown
// documents.
const FileExtension = "*.au.md"

// Loader loads [test.Node] values from (directories containing) Markdown documents.
type Loader struct {
	FS     fs.FS
	Parser parser.Parser
}

// Load returns a [Test] comprising the Aureus Markdown documents at the the
// given filesystem path.
//
// name may be the path to a single document or a directory containing multiple
// documents. If name is a directory, all documents with an extensioning
// matching [FileExtension] are loaded recursively.
//
// Any directory or file that begins with an underscore is marked as skipped.
func (l *Loader) Load(name string) (test.Test, error) {
	info, err := fs.Stat(l.filesystem(), name)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return l.fromDir(name, test.FlagNone)
	}

	return l.fromFile(name, test.FlagNone)
}

func (l *Loader) fromDir(dir string, flags test.Flags) (test.Test, error) {
	entries, err := fs.ReadDir(l.filesystem(), dir)
	if err != nil {
		return nil, err
	}

	name := path.Base(dir)
	name, skip := strings.CutPrefix(name, "_")
	suite := test.Suite{
		Name:   name,
		Flags:  flags,
		Origin: test.Directory{Name: dir},
	}

	if skip {
		suite.Flags.Add(test.FlagSkipped)
		flags.Add(test.FlagAncestorSkipped)
	}

	for _, entry := range entries {
		n := path.Join(dir, entry.Name())

		load := l.fromFile
		if entry.IsDir() {
			load = l.fromDir
		} else if !strings.HasSuffix(entry.Name(), FileExtension) {
			continue
		}

		sub, err := load(n, flags)
		if err != nil {
			return nil, err
		}

		suite.Tests = append(suite.Tests, sub)
	}

	return suite, nil
}

func (l *Loader) fromFile(filename string, flags test.Flags) (test.Test, error) {
	name := path.Base(filename)
	name, skip := strings.CutPrefix(name, "_")
	suite := test.Suite{
		Name:   name,
		Flags:  flags,
		Origin: test.Document{Name: filename},
	}

	if skip {
		suite.Flags.Add(test.FlagSkipped)
		flags.Add(test.FlagAncestorSkipped)
	}

	source, err := fs.ReadFile(l.filesystem(), filename)
	if err != nil {
		return nil, err
	}

	root := l.parser().Parse(
		text.NewReader(source),
	)

	if err := l.fromNode(&suite, root); err != nil {
		return nil, err
	}

	return suite, nil
}

func (l *Loader) fromNode(suite *test.Suite, doc ast.Node) error {
	// fmt.Println(
	// 	strings.Repeat("  ", depth),
	// 	"-",
	// 	node.Kind(),
	// )

	// var blocks []test.CodeBlock

	level := 0
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		switch n := n.(type) {
		case *ast.Heading:
			// if n.Level > level {
			// 	level = n.Level
			// }
		case *ast.FencedCodeBlock:

		}
	}

	// 	}
	// 	status, err = walker(n, false)
	// 	if err != nil || status == WalkStop {
	// 		return WalkStop, err
	// 	}
	// 	return WalkContinue, nil
	// }

	// var assertion *EqualAssertion

	// return test, ast.Walk(
	// 	root,
	// 	func(n ast.Node, entering bool) (ast.WalkStatus, error) {
	// 		if !entering {
	// 			if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
	// 				test.Name = string(h.Text(source))
	// 			}
	// 			return ast.WalkContinue, nil
	// 		}

	// 		code, ok := n.(*ast.FencedCodeBlock)
	// 		if !ok {
	// 			return ast.WalkContinue, nil
	// 		}

	// 		line := astx.LineNumber(code, source)

	// 		info, err := l.parseCodeInfo(code, source)
	// 		if err != nil {
	// 			return ast.WalkStop, fmt.Errorf(
	// 				"aureus assertion at %s:%d is invalid: %s",
	// 				filename,
	// 				line,
	// 				err,
	// 			)
	// 		}

	// 		if !info.IsAssertion {
	// 			assertion = &EqualAssertion{
	// 				Name:          "L" + strconv.Itoa(line),
	// 				File:          filename,
	// 				Line:          line,
	// 				InputLanguage: info.Language,
	// 				Input:         astx.Source(code, source),
	// 			}
	// 			return ast.WalkContinue, nil
	// 		} else if assertion == nil {
	// 			return ast.WalkStop, fmt.Errorf(
	// 				"assertion in %s on line %d: preceding code block containing the input value was not found",
	// 				filename,
	// 				line,
	// 			)
	// 		} else {
	// 			assertion.OutputLanguage = info.Language
	// 			assertion.ExpectedOutput = astx.Source(code, source)
	// 			assertion.Skip = info.Skip
	// 			test.Assertions = append(test.Assertions, *assertion)
	// 			assertion = nil
	// 		}

	// 		return ast.WalkContinue, nil
	// 	},
	// )
	return nil
}

func (l *Loader) filesystem() fs.FS {
	if l.FS == nil {
		return rootFS{}
	}
	return l.FS
}

func (l *Loader) parser() parser.Parser {
	if l.Parser == nil {
		return goldmark.DefaultParser()
	}
	return l.Parser
}
