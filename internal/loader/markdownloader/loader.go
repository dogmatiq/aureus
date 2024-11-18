package markdownloader

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dogmatiq/aureus/internal/loader"
	"github.com/dogmatiq/aureus/internal/rootfs"
	"github.com/dogmatiq/aureus/internal/test"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Loader loads [test.Test] values from Markdown files containing code blocks
// representing test inputs and expected outputs.
type Loader struct {
	options loadOptions
}

// NewLoader returns a new [Loader], which loads golden file tests from the
// filesystem.
func NewLoader(options ...LoadOption) *Loader {
	l := &Loader{
		options: loadOptions{
			FS:          rootfs.FS,
			Recurse:     true,
			LoadContent: LoadContent,
			Parser:      goldmark.DefaultParser(),
		},
	}

	for _, opt := range options {
		opt(&l.options)
	}

	return l
}

// Load returns a test built from files in the given directory.
//
// Any directory or file that begins with an underscore produces a test that is
// marked as skipped.
func (l *Loader) Load(dir string, options ...LoadOption) (test.Test, error) {
	opts := l.options
	for _, opt := range options {
		opt(&opts)
	}

	return loader.LoadDir(
		opts.FS,
		dir,
		opts.Recurse,
		func(builder *loader.TestBuilder, fsys fs.FS, filePath string) error {
			return loadFile(builder, opts, filePath)
		},
	)
}

func loadFile(builder *loader.TestBuilder, opts loadOptions, filePath string) error {
	name := path.Base(filePath)
	name, skip := strings.CutPrefix(name, "_")
	name, ok := strings.CutSuffix(name, ".md")
	if !ok {
		return nil
	}

	source, err := fs.ReadFile(opts.FS, filePath)
	if err != nil {
		return err
	}

	title, tests, err := loadDocument(
		opts,
		filePath,
		source,
	)
	if err != nil {
		return err
	}

	if title != "" {
		name = title
	}

	builder.AddTest(
		test.New(
			name,
			test.WithSkip(skip),
			test.WithSubTests(tests...),
		),
	)

	return nil
}

func loadDocument(
	opts loadOptions,
	filePath string,
	source []byte,
) (string, []test.Test, error) {
	doc := opts.Parser.Parse(
		text.NewReader(source),
	)

	var (
		b        loader.TestBuilder
		title    string
		headings []string
	)

	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		switch n := n.(type) {
		case *ast.Heading:
			if n.Level == 1 {
				if n == doc.FirstChild() {
					// If we have a level-one heading that's the first node in
					// the document treat it as the title of the document.
					title = linesOf(n, source)
				} else {
					// Unless there are other level-one headings later in the
					// document.
					title = ""
				}
			}

			for len(headings) < n.Level {
				headings = append(headings, "")
			}
			headings[n.Level-1] = linesOf(n, source)
			headings = headings[:n.Level]

		case *ast.FencedCodeBlock:
			if err := loadBlock(
				&b,
				opts,
				filePath,
				source,
				headings,
				n,
			); err != nil {
				return "", nil, err
			}
		}
	}

	tests, err := b.Build()
	return title, tests, err
}

func loadBlock(
	builder *loader.TestBuilder,
	opts loadOptions,
	filePath string,
	source []byte,
	headings []string,
	block *ast.FencedCodeBlock,
) error {
	info := ""
	if block.Info != nil {
		info = string(block.Info.Value(source))
	}

	code := linesOf(block, source)

	content, skip, err := opts.LoadContent(headings, info, code)
	if err != nil {
		return err
	}

	line, begin, end := locationOf(block, source)

	return builder.AddContent(
		loader.ContentEnvelope{
			File:    filePath,
			Line:    line,
			Begin:   int64(begin),
			End:     int64(end),
			Skip:    skip,
			Content: content,
		},
	)
}
