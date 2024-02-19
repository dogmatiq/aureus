package markdownloader

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dogmatiq/aureus/internal/loader"
	"golang.org/x/net/html"
)

// A ContentLoader is a function that returns content obtained from a fenced
// code block.
//
// info is the "info string" of the code block, that is the line of text that
// contains the language and other loader-specific information.
//
// See https://spec.commonmark.org/0.31.2/#fenced-code-blocks
type ContentLoader func(info, code string) (_ loader.Content, skip bool, _ error)

// LoadContent is the default [ContentLoader] implementation.
func LoadContent(info, code string) (loader.Content, bool, error) {
	lang, attrs, err := ParseInfoString(info)
	if err != nil {
		return loader.Content{}, false, err
	}

	isInput, err := extractFlag(attrs, inputAttr)
	if err != nil {
		return loader.Content{}, false, err
	}

	isOutput, err := extractFlag(attrs, outputAttr)
	if err != nil {
		return loader.Content{}, false, err
	}

	if isInput && isOutput {
		return loader.Content{}, false, fmt.Errorf(
			"only one of '%s%s' and '%s%s' may be specified",
			attrPrefix, inputAttr,
			attrPrefix, outputAttr,
		)
	}

	group, err := extractValue(attrs, groupAttr)
	if err != nil {
		return loader.Content{}, false, err
	}

	skip, err := extractFlag(attrs, skipAttr)
	if err != nil {
		return loader.Content{}, false, err
	}

	for k := range attrs {
		if strings.HasPrefix(k, attrPrefix) {
			return loader.Content{}, false, fmt.Errorf("unrecognized attribute %q", k)
		}
	}

	c := loader.Content{
		Group:      group,
		Language:   lang,
		Attributes: attrs,
		Data:       code,
	}

	if isInput {
		c.Role = loader.Input
	} else if isOutput {
		c.Role = loader.Output
	} else {
		return loader.Content{}, false, nil
	}

	return c, skip, nil
}

// ParseInfoString parses the "info string" of a fenced code block into a
// language and a map of attributes.
func ParseInfoString(info string) (lang string, attrs map[string]string, err error) {
	data := &bytes.Buffer{}
	data.WriteString("<html ")
	data.WriteString(info)
	data.WriteByte('>')

	node, err := html.Parse(data)
	if err != nil {
		return "", nil, err
	}

	attrs = map[string]string{}
	for i, attr := range node.FirstChild.Attr {
		if i == 0 && attr.Val == "" && !strings.HasPrefix(attr.Key, attrPrefix) {
			lang = attr.Key
		} else {
			attrs[attr.Key] = attr.Val
		}
	}

	return lang, attrs, nil
}

const (
	attrPrefix = "au:"
	inputAttr  = "input"
	outputAttr = "output"
	groupAttr  = "group"
	skipAttr   = "skip"
)

func extractFlag(attrs map[string]string, k string) (bool, error) {
	k = attrPrefix + k
	v, ok := attrs[k]
	if v != "" {
		return false, fmt.Errorf("%q attribute must not have a value", k)
	}
	delete(attrs, k)
	return ok, nil
}

func extractValue(attrs map[string]string, k string) (string, error) {
	k = attrPrefix + k
	v, ok := attrs[k]
	if !ok {
		return "", nil
	}
	if v == "" {
		return "", fmt.Errorf("%q attribute must have a value", k)
	}
	delete(attrs, k)
	return v, nil
}
