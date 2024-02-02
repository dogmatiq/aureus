package test

// Assertion is an interface for an assertion made by a [Test].
type Assertion interface {
	AcceptVisitor(AssertionVisitor)
}

// AssertionVisitor is an interface for dispatching based on the concrete type
// of an [Assertion].
type AssertionVisitor interface {
	VisitEqualAssertion(*EqualAssertion)
	VisitDiffAssertion(*DiffAssertion)
}

// EqualAssertion is an [Assertion] that asserts that asserts two values are equal.
type EqualAssertion struct {
	Input  Content
	Output Content
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (a *EqualAssertion) AcceptVisitor(v AssertionVisitor) {
	v.VisitEqualAssertion(a)
}

// DiffAssertion is an [Assertion] that asserts that two values differ in a specific way.
type DiffAssertion struct {
	Input Content
	Diff  Content
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (a *DiffAssertion) AcceptVisitor(v AssertionVisitor) {
	v.VisitDiffAssertion(a)
}
