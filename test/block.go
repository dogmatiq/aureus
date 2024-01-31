package test

// CodeBlock contains information about a fenced code block, used to specify inputs
// and outputs for tests.
type CodeBlock struct {
	Origin     Origin
	Language   string
	Text       string
	Attributes map[string]string
}
