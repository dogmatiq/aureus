package loader

import (
	"github.com/dogmatiq/aureus/internal/test"
)

// ContentRole is an enumeration of the roles that loaded content can play
// within in a test.
type ContentRole int

const (
	// NoRole indicates that the content has no specific role within the test.
	NoRole ContentRole = iota

	// Input indicates that the content is an input to a test.
	Input

	// Output indicates that the content is the expected output from a test.
	Output
)

// Content is a specialization of [test.Content] that includes meta-data about
// how it was loaded and how it should appear within tests.
type Content struct {
	// Role is the role that the content plays within a test. A value of
	// [NoRole] indicates that the content has no specific role within the test.
	Role ContentRole

	// Group is the name of the group to which the content belongs. Inputs and
	// outputs in the same group form a matrix of test cases.
	Group *Group

	// Caption is an optional disambiguating name, title or short description of
	// the content.
	Caption string

	// Language is the language of the content, if known, e.g. "json", "yaml",
	// etc. Content with an empty language is treated as plain text.
	Language string

	// Attributes is a set of key-value pairs that provide additional
	// loader-specific information about the data.
	Attributes map[string]string

	// Data is the content itself.
	Data []byte
}

// ContentEnvelope is a container for [Content] and meta-data about how it was
// loaded.
type ContentEnvelope struct {
	// File is the path of the file from which the content was loaded.
	File string

	// Line is the line number within the file where the content begins, or 0 if
	// the content represents the entire file.
	Line int

	// The half-open range [Begin, End) is the section within the file that
	// contains the content, given in bytes.
	//
	// If the range is [0, 0), the content represents the entire file.
	Begin, End int64

	// Skip is a flag that indicates whether this content should be skipped when
	// running tests.
	Skip bool

	// Content is the loaded content.
	Content Content
}

// AsTestContent returns the content as a [test.Content].
func (e ContentEnvelope) AsTestContent() test.Content {
	return test.Content{
		ContentMetaData: test.ContentMetaData{
			File:       e.File,
			Line:       e.Line,
			Begin:      e.Begin,
			End:        e.End,
			Language:   e.Content.Language,
			Attributes: e.Content.Attributes,
		},
		Data: e.Content.Data,
	}
}

// IsEntireFile returns true if the content occupies the entire file.
func (e ContentEnvelope) IsEntireFile() bool {
	return e.Begin == 0 && e.End == 0
}

// SeparateContentByRole separates content into inputs and outputs.
func SeparateContentByRole(content []ContentEnvelope) (inputs, outputs []ContentEnvelope) {
	for _, c := range content {
		switch c.Content.Role {
		case Input:
			inputs = append(inputs, c)
		case Output:
			outputs = append(outputs, c)
		}
	}

	return inputs, outputs
}

// Group is the group to which content belongs. Inputs and outputs in the same
// group form a matrix of test cases.
type Group struct {
	name string
}

// IsNamed returns true if the group has a name.
func (g *Group) IsNamed() bool {
	return g.name != ""
}

// Name returns the name of the group.
//
// It panics if the group is unnamed.
func (g *Group) Name() string {
	if !g.IsNamed() {
		panic("group is un-named")
	}
	return g.name
}

// NamedGroup returns a [Group] with the given name.
func NamedGroup(name string) *Group {
	if name == "" {
		panic("group name must not be empty")
	}
	return &Group{name}
}

// UnnamedGroup returns a [Group] with no name.
func UnnamedGroup() *Group {
	return &Group{}
}

// groupName returns the name of g as a string.
func groupName(g *Group) string {
	if g == nil {
		return ""
	}
	return g.name
}
