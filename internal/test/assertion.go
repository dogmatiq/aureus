package test

// Assertion is an interface for an assertion made by a [Test].
type Assertion interface {
	AcceptVisitor(AssertionVisitor, ...VisitOption)
}

// AssertionVisitor is an interface for dispatching based on the concrete type
// of an [Assertion].
type AssertionVisitor interface {
	VisitEqualAssertion(EqualAssertion)
}

// EqualAssertion is an [Assertion] that asserts that asserts two values are
// equal.
type EqualAssertion struct {
	Input  Content
	Output Content
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of a.
func (a EqualAssertion) AcceptVisitor(v AssertionVisitor, options ...VisitOption) {
	cfg := newVisitConfig(options)
	if cfg.TestingT != nil {
		cfg.TestingT.Helper()
	}
	v.VisitEqualAssertion(a)
}