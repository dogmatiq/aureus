package test

// Content contains input or output content.
type Content struct {
	Origin   Origin
	MetaData map[string]string
	Data     []byte
}
