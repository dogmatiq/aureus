package test

// TestingT is the subset of the [*testing.T] interface that is used when
// dispatching a visitor to ensure that the internals of the visitor mechanism
// are properly marked as test helpers.
type TestingT interface {
	Helper()
}

// WithT configures the visitor to use t to mark the internal mechanism of the
// visitor as test helpers.
func WithT(t TestingT) VisitOption {
	return func(c *visitConfig) {
		c.TestingT = t
	}
}

// VisitOption is an option for configuring a visitor.
type VisitOption func(*visitConfig)

type visitConfig struct {
	TestingT TestingT
}

func newVisitConfig(opts []VisitOption) visitConfig {
	var cfg visitConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
