package runner

// LoggerT is the subset of the [testing.TB] that supports logging only.
type LoggerT interface {
	Helper()
	Name() string
	Log(...any)
}

// FailerT is the subset of the [testing.TB] that supports logging and failure
// reporting.
type FailerT interface {
	LoggerT
	SkipNow()
	Fail()
	Failed() bool
}

// TestingT is a constraint for types that are compatible with [testing.T].
type TestingT[T any] interface {
	FailerT

	SkipNow()
	Fail()
	Failed() bool
	Run(string, func(T)) bool
}
