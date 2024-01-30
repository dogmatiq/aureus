package aureus

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const (
	defaultFileExtension = ".au.md"
)

// Loader loads [Test] values from (directories containing) Markdown documents.
type Loader struct {
	FS              fs.FS
	Parser          parser.Parser
	FileExtension   string
	DefaultLanguage string
}

// Load returns tests loaded from the file or directory at the given path.
func (l *Loader) Load(ctx context.Context, p string) (Test, error) {
	if ctx.Err() != nil {
		return Suite{}, ctx.Err()
	}

	info, err := fs.Stat(l.FS, p)
	if err != nil {
		return nil, err
	}

	dir := path.Dir(p)

	if info.IsDir() {
		return l.loadSuite(ctx, dir, info.Name())
	}

	return l.loadDocument(ctx, dir, info.Name())
}

func (l *Loader) loadSuite(ctx context.Context, dir, name string) (Suite, error) {
	if ctx.Err() != nil {
		return Suite{}, ctx.Err()
	}

	qual := path.Join(dir, name)
	name, skip := strings.CutPrefix(name, "_")

	entries, err := fs.ReadDir(l.FS, qual)
	if err != nil {
		return Suite{}, err
	}

	test := Suite{
		Name: name,
		Dir:  qual,
		Skip: skip,
	}

	for _, entry := range entries {
		var sub Test

		if entry.IsDir() {
			sub, err = l.loadSuite(ctx, qual, entry.Name())
		} else if l.isTestDocument(entry.Name()) {
			sub, err = l.loadDocument(ctx, qual, entry.Name())
		} else {
			continue
		}

		if err != nil {
			return Suite{}, err
		}

		test.Tests = append(test.Tests, sub)
	}

	return test, nil
}

func (l *Loader) loadDocument(ctx context.Context, dir, name string) (Document, error) {
	if ctx.Err() != nil {
		return Document{}, ctx.Err()
	}

	qual := path.Join(dir, name)
	name, skip := strings.CutPrefix(name, "_")

	if !l.isTestDocument(qual) {
		return Document{}, fmt.Errorf("%s is not an aureus document", qual)
	}

	source, err := fs.ReadFile(l.FS, qual)
	if err != nil {
		return Document{}, err
	}

	parser := l.Parser
	if parser == nil {
		parser = goldmark.DefaultParser()
	}
	root := parser.Parse(text.NewReader(source))

	test := Document{
		Name: name,
		File: qual,
		Skip: skip,
	}

	var assertion *Assertion

	return test, ast.Walk(
		root,
		func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
					test.Name = string(h.Text(source))
				}
				return ast.WalkContinue, nil
			}

			code, ok := n.(*ast.FencedCodeBlock)
			if !ok {
				return ast.WalkContinue, nil
			}

			info := l.parseCodeInfo(code, source)
			line := lineNumberOf(code, source)

			if !info.IsAssertion {
				assertion = &Assertion{
					Name:          "L" + strconv.Itoa(line),
					File:          qual,
					Line:          line,
					InputLanguage: info.Language,
					Input:         sourceOf(code, source),
				}
				return ast.WalkContinue, nil
			} else if assertion == nil {
				test.Errors = append(
					test.Errors,
					fmt.Errorf(
						"aureus assertion at %s:%d does not have a preceding code block containing the input value",
						qual,
						line,
					),
				)
			} else {
				assertion.OutputLanguage = info.Language
				assertion.ExpectedOutput = sourceOf(code, source)
				assertion.Skip = info.Skip
				test.Assertions = append(test.Assertions, *assertion)
				assertion = nil
			}

			return ast.WalkContinue, nil
		},
	)
}

// isTestDocument returns true if the given file should be treated as a Markdown
// document that (potentially) contains test assertions.
func (l *Loader) isTestDocument(filename string) bool {
	ext := l.FileExtension
	if ext == "" {
		ext = defaultFileExtension
	}
	return strings.HasSuffix(filename, ext)
}

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
