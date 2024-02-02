package test

// Runnable is an interface for a test.
type Runnable interface {
	AcceptVisitor(Visitor)
}

// Visitor is an interface for dispatching based on the concrete type of a
// [Runnable].
type Visitor interface {
	VisitSuite(*Suite)
	VisitTest(*Test)
}

// Suite is a [Runnable] that contains a collection of sub-tests.
type Suite struct {
	Name   string
	Flags  FlagSet
	Origin Origin
	Tests  []Runnable `json:",omitempty"`
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (n *Suite) AcceptVisitor(v Visitor) {
	v.VisitSuite(n)
}

// Test is a [Runnable] that makes one or more assertions.
type Test struct {
	Name       string
	Flags      FlagSet
	Origin     Origin
	Assertions []Assertion `json:",omitempty"`
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of n.
func (n *Test) AcceptVisitor(v Visitor) {
	v.VisitTest(n)
}
