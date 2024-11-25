package runner

// assertionExecutor is an impelmentation of [test.AssertionVisitor] that
// performs assertions within the context of a test.
type assertionExecutor[T TestingT[T]] struct {
	TestingT T
	Runner   *Runner[T]
}
