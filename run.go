package aureus

import (
	"io"

	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	"github.com/dogmatiq/aureus/internal/loader/markdownloader"
	"github.com/dogmatiq/aureus/internal/runner"
	"github.com/dogmatiq/aureus/internal/test"
)

// OutputGenerator produced output for a test case's input.
//
// in is the input data for the test case. out is meta-data about the expected
// output, which may be used to influence the kind of output that the generator
// writes to w.
type OutputGenerator func(
	w io.Writer,
	in Content,
	out ContentMetaData,
) error

type (
	// Content is data used as input or output in tests.
	Content = test.Content
	// ContentMetaData contains information about input or output content.
	ContentMetaData = test.ContentMetaData
)

// Run searches a directory for "golden" tests and executes them as
// sub-tests of t.
//
// By default is searches the ./testdata directory for test cases. This can be
// changed using the [WithDir] option.
//
// g is a function that generates output from input values. It is called for
// each test case. If the output it produces does not match the expected output.
func Run[T runner.TestingT[T]](t T, g OutputGenerator, options ...RunOption) {
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
		GenerateOutput: (runner.OutputGenerator)(g),
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
