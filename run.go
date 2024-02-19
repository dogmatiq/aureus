package aureus

import (
	"github.com/dogmatiq/aureus/loader/fileloader"
	"github.com/dogmatiq/aureus/loader/markdownloader"
	"github.com/dogmatiq/aureus/runner"
)

// Run searches a directory for "golden" tests and executes them as
// sub-tests of t.
//
// By default is searches the ./testdata directory for test cases. This can be
// changed using the [WithDir] option.
//
// g is a function that generates output from input values. It is called for
// each test case. If the output it produces does not match the expected output.
func Run[T runner.TestingT[T]](t T, g runner.OutputGenerator, options ...RunOption) {
	t.Helper()

	opts := runOptions{
		Dir:       "./testdata",
		Recursive: true,
	}
	for _, opt := range options {
		opt(&opts)
	}

	fileLoader := fileloader.NewLoader(fileloader.WithRecursion(opts.Recursive))
	fileTests, err := fileLoader.Load(opts.Dir)
	if err != nil {
		t.Log("failed to load test cases:", err)
		t.Fail()
		return
	}

	markdownLoader := markdownloader.NewLoader(markdownloader.WithRecursion(opts.Recursive))
	markdownTests, err := markdownLoader.Load(opts.Dir)
	if err != nil {
		t.Log("failed to load test cases:", err)
		t.Fail()
		return
	}

	r := runner.Runner[T]{
		GenerateOutput: g,
	}
	r.Run(t, fileTests)
	r.Run(t, markdownTests)
}

// RunOption is an option that changes the behavior of [Run].
type RunOption func(*runOptions)

type runOptions struct {
	Dir       string
	Recursive bool
}

// WithDir is a [RunOption] that sets the directory to search for test cases.
// By default the ./testdata directory is used.
func WithDir(dir string) RunOption {
	return func(o *runOptions) {
		o.Dir = dir
	}
}

// WithRecursion is a [RunOption] that enables or disables recursion when
// searching for test cases. By default recursion is enabled.
func WithRecursion(on bool) RunOption {
	return func(o *runOptions) {
		o.Recursive = on
	}
}
