package runner

// TestingT is a constraint for a type that is compatible with [testing.T].
type TestingT[T any] interface {
	Helper()
	Parallel()
	Run(string, func(T)) bool
	Log(...any)
	SkipNow()
	Fail()
}
