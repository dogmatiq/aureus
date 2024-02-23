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

	var docBuilder loader.TestBuilder
	doc := opts.Parser.Parse(
		text.NewReader(source),
	)

	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if block, ok := node.(*ast.FencedCodeBlock); ok {
			if err := loadBlock(
				&docBuilder,
				opts,
				filePath,
				source,
				block,
			); err != nil {
				return err
			}
		}
	}

	subTests, err := docBuilder.Build()
	if err != nil {
		return err
	}

	builder.AddTest(
		test.New(
			name,
			test.WithSkip(skip),
			test.WithSubTests(subTests...),
		),
	)

	return nil
}

func loadBlock(
	builder *loader.TestBuilder,
	opts loadOptions,
	filePath string,
	source []byte,
	block *ast.FencedCodeBlock,
) error {
	info := ""
	if block.Info != nil {
		info = string(block.Info.Text(source))
	}
	code := sourceOf(block, source)

	c, skip, err := opts.LoadContent(info, code)
	if err != nil {
		return err
	}

	return builder.AddContent(
		loader.ContentEnvelope{
			File:    filePath,
			Line:    lineNumberOf(block, source),
			Skip:    skip,
			Content: c,
		},
	)
}
