package test

// Assertion is an interface for an assertion made by a [Test].
type Assertion interface {
	AcceptVisitor(AssertionVisitor)
}

// AssertionVisitor is an interface for dispatching based on the concrete type
// of an [Assertion].
type AssertionVisitor interface {
	VisitEqual(*Equal)
	VisitDiff(*Diff)
}

// Equal is an [Assertion] that asserts that asserts two values are equal.
type Equal struct {
	Input  Content
	Output Content
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (a *Equal) AcceptVisitor(v AssertionVisitor) {
	v.VisitEqual(a)
}

// Diff is an [Assertion] that asserts that two values differ in a specific way.
type Diff struct {
	Input Content
	Diff  Content
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (a *Diff) AcceptVisitor(v AssertionVisitor) {
	v.VisitDiff(a)
}
