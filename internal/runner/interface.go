package runner

// TestingT is a constraint for types that are compatible with [testing.T].
type TestingT[T any] interface {
	Helper()
	Log(...any)
	SkipNow()
	Fail()
	Failed() bool
	Run(string, func(T)) bool
}
