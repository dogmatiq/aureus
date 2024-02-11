package test

// Content contains input or output content.
type Content struct {
	File       string
	Line       int `json:",omitempty"` // 0 == whole file
	Data       []byte
	Attributes map[string]string `json:",omitempty"`
}
