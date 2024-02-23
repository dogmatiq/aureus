package aureus

import (
	"bytes"
	"io"

	"github.com/dogmatiq/aureus/internal/loader/fileloader"
	"github.com/dogmatiq/aureus/internal/loader/markdownloader"
	"github.com/dogmatiq/aureus/internal/runner"
	"github.com/dogmatiq/aureus/internal/test"
)

// TestingT is a constraint for the subset of [testing.T] that is used by Aureus
// to execute tests.
type TestingT[T any] interface {
	Helper()
	Run(string, func(T)) bool
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
func Run[T runner.TestingT[T]](t T, g OutputGenerator[T], options ...RunOption) {
	t.Helper()

	opts := runOptions{
		Dir:       "./testdata",
		Recursive: true,
		TrimSpace: true,
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
		GenerateOutput: func(
			t T,
			w io.Writer,
			in test.Content,
			out test.ContentMetaData,
		) error {
			return g(
				t,
				&input{
					Reader: bytes.NewReader(in.Data),
					lang:   in.Language,
					attrs:  in.Attributes,
				},
				&output{
					Writer: w,
					lang:   out.Language,
					attrs:  out.Attributes,
				},
			)
		},
		TrimSpace: opts.TrimSpace,
	}
	r.Run(t, fileTests)
	r.Run(t, markdownTests)
}

// RunOption is an option that changes the behavior of [Run].
type RunOption func(*runOptions)

type runOptions struct {
	Dir       string
	Recursive bool
	TrimSpace bool
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
