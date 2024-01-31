package test

// Test is an interface for a test.
type Test interface {
	AcceptVisitor(Visitor)
}

// Visitor is an interface for dispatching based on the concrete type of a
// [Test].
type Visitor interface {
	VisitSuite(Suite)
	VisitEqual(Equal)
	VisitDiff(Diff)
}

// Suite is a [Test] that contains a collection of sub-tests.
type Suite struct {
	Name   string
	Flags  Flags
	Origin Origin
	Tests  []Test
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (n Suite) AcceptVisitor(v Visitor) {
	v.VisitSuite(n)
}

// Equal is a [Test] that asserts that its output is equal to a specific value.
type Equal struct {
	Name   string
	Flags  Flags
	Input  CodeBlock
	Expect CodeBlock
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (n Equal) AcceptVisitor(v Visitor) {
	v.VisitEqual(n)
}

// Diff is a [Test] that asserts that its output is different to a specific
// value, specified as a diff.
type Diff struct {
	Name   string
	Flags  Flags
	Input  CodeBlock
	Expect CodeBlock
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (n Diff) AcceptVisitor(v Visitor) {
	v.VisitDiff(n)
}
