package test

// Test is a (possibly nested) test.
type Test struct {
	Name       string
	Skip       bool
	SubTests   []Test
	Assertions []Assertion
}

// Assertion represents a requirement that an input match a specific output.
type Assertion struct {
	Input  Content
	Output Content
}

// New creates a new [Test].
//
// It returns the test and the set of flags that should be inherited by any
// sub-tests.
func New(name string, options ...Option) Test {
	t := Test{
		Name: name,
	}

	for _, opt := range options {
		opt(&t)
	}

	return t
}

// Option is an option that controls how a test is created by [New].
type Option func(*Test)

// WithSkip is a [TestOption] that sets the skip flag.
func WithSkip(skip bool) Option {
	return func(t *Test) {
		t.Skip = skip
	}
}

// WithSubTests is a [TestOption] that adds sub-tests to the test.
func WithSubTests(subTests ...Test) Option {
	return func(t *Test) {
		t.SubTests = append(t.SubTests, subTests...)
	}
}

// WithAssertions is a [TestOption] thatadds assertions to the test.
func WithAssertions(assertions ...Assertion) Option {
	return func(opts *Test) {
		opts.Assertions = append(opts.Assertions, assertions...)
	}
}
