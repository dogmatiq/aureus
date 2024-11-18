package aureus

import (
	"github.com/dogmatiq/aureus/internal/cliflags"
	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	"github.com/dogmatiq/aureus/internal/loader/markdownloader"
	"github.com/dogmatiq/aureus/internal/runner"
	"github.com/dogmatiq/aureus/internal/test"
)

// TestingT is a constraint for the subset of [testing.T] that is used by Aureus
// to execute tests.
type TestingT[Self any] interface {
	Helper()
	Run(string, func(Self)) bool
	Log(...any)
	SkipNow()
	Fail()
}

// Run searches a directory for tests and executes them as sub-tests of t.
//
// By default it searches the ./testdata directory; use the [FromDir] option to
// search a different directory.
//
// g is an [OutputGenerator] that produces output from input values for each
// test. If the output produced by g does not match the test's expected output
// the test fails.
func Run[T runner.TestingT[T]](
	t T,
	g OutputGenerator[T],
	options ...RunOption,
) {
	t.Helper()

	opts := runOptions{
		Dir:           "./testdata",
		Recursive:     true,
		TrimSpace:     true,
		BlessStrategy: &runner.BlessAvailable{},
	}

	if cliflags.Get().Bless {
		Bless(true)(&opts)
	}

	for _, opt := range options {
		opt(&opts)
	}

	fileLoader := fileloader.NewLoader(fileloader.WithRecursion(opts.Recursive))
	fileTests, err := fileLoader.Load(opts.Dir)
	if err != nil {
		t.Log("failed to load tests:", err)
		t.Fail()
		return
	}

	markdownLoader := markdownloader.NewLoader(markdownloader.WithRecursion(opts.Recursive))
	markdownTests, err := markdownLoader.Load(opts.Dir)
	if err != nil {
		t.Log("failed to load tests:", err)
		t.Fail()
		return
	}

	r := runner.Runner[T]{
		GenerateOutput: func(t T, in runner.Input, out runner.Output) error {
			return g(t, in, out)
		},
		TrimSpace:     opts.TrimSpace,
		BlessStrategy: opts.BlessStrategy,
	}

	tests := test.Merge(fileTests, markdownTests)

	if len(tests) == 0 {
		t.Log("no tests found")
	} else {
		for _, x := range tests {
			r.Run(t, x)
		}
	}
}

// RunOption is an option that changes the behavior of [Run].
type RunOption func(*runOptions)

type runOptions struct {
	Dir           string
	Recursive     bool
	TrimSpace     bool
	BlessStrategy runner.BlessStrategy
}

// FromDir is a [RunOption] that sets the directory to search for tests. By
// default the ./testdata directory is used.
func FromDir(dir string) RunOption {
	return func(o *runOptions) {
		o.Dir = dir
	}
}

// Recursive is a [RunOption] that enables or disables recursion when searching
// for test cases. By default recursion is enabled.
func Recursive(on bool) RunOption {
	return func(o *runOptions) {
		o.Recursive = on
	}
}

// TrimSpace is a [RunOption] that enables or disables trimming of leading and
// trailing whitespace from test outputs. By default trimming is enabled.
func TrimSpace(on bool) RunOption {
	return func(o *runOptions) {
		o.TrimSpace = on
	}
}

// Bless is a [RunOption] that enables or disables "blessing" of failed tests.
//
// If blessing is enabled, the file containing the expected output of each
// failed assertion is replaced with the actual output.
//
// By default blessing is disabled unless the -aureus.bless flag is set on the
// command line.
func Bless(on bool) RunOption {
	return func(o *runOptions) {
		if on {
			o.BlessStrategy = &runner.BlessEnabled{}
		} else {
			o.BlessStrategy = &runner.BlessDisabled{}
		}
	}
}
