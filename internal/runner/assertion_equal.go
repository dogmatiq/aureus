package runner

import (
	"os"

	"github.com/dogmatiq/aureus/internal/test"
)

func (x *assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	got, err := x.generateOutput(a.Input, a.Output)
	if err != nil {
		x.TestingT.Log(err)
		x.TestingT.Fail()
		return
	}
	defer func() {
		got.Close()
		if !x.TestingT.Failed() {
			os.Remove(got.Name())
		}
	}()

	want, err := a.Output.Open()
	if err != nil {
		x.TestingT.Log("unable to open output file:", err)
		x.TestingT.Fail()
		return
	}
	defer want.Close()

	diff, err := x.computeDiff(
		location(a.Output), want,
		got.Name(), got,
	)
	if err != nil {
		x.TestingT.Log("unable to generate diff:", err)
		x.TestingT.Fail()
		return
	}

	if diff != nil {
		x.logSection("DIFF", diff)
		x.TestingT.Fail()
	}
}
