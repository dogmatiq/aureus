package runner

// LoggerT is the subset of the [testing.TB] that supports logging only.
type LoggerT interface {
	Helper()
	Name() string
	Log(...any)
}

// TestingT is a constraint for types that are compatible with [testing.T].
type TestingT[T any] interface {
	LoggerT

	SkipNow()
	Fail()
	Failed() bool
	Run(string, func(T)) bool
}
