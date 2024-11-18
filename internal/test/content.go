package test

// Content is data used as input or output in tests.
type Content struct {
	// ContentMetaData is additional information about the content.
	ContentMetaData

	// Data is the content itself.
	Data []byte
}

// ContentMetaData contains information about input or output content.
type ContentMetaData struct {
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

	// Language is the language of the content, if known, e.g. "json", "yaml",
	// etc. Content with an empty language is treated as plain text.
	Language string

	// Attributes is a set of key-value pairs that provide additional
	// loader-specific information about the data.
	Attributes map[string]string
}

// IsEntireFile returns true if the content occupies the entire file.
func (m ContentMetaData) IsEntireFile() bool {
	return m.Begin == 0 && m.End == 0
}
