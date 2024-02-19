package test

// Content is data used as input or output in tests.
type Content struct {
	// ContentMetaData is additional information about the content.
	ContentMetaData

	// Data is the content itself. It is always considered to be human-readable
	// text, such as source code.
	Data string
}

// ContentMetaData is additional information about a piece of content.
type ContentMetaData struct {
	// File is the path of the file from which the content was loaded.
	File string

	// Line is the line number within the file where the content begins, or 0 if
	// the content represents the entire file.
	Line int `json:",omitempty"`
	// Language is the language of the content, if known, e.g. "json", "yaml",
	// etc. Content with an empty language is treated as plain text.
	Language string `json:",omitempty"`

	// Attributes is a set of key-value pairs that provide additional
	// loader-specific information about the data.
	Attributes map[string]string `json:",omitempty"`
}
