package runner

import (
	"bytes"
	"os"

	"github.com/dogmatiq/aureus/internal/test"
)

func (x *assertionExecutor[T]) VisitEqualAssertion(a test.EqualAssertion) {
	x.TestingT.Helper()

	x.logSection("INPUT", a.Input.Data, location(a.Input))

	f, err := x.generateOutput(a.Input, a.Output)
	if err != nil {
		x.TestingT.Log(err)
		x.TestingT.Fail()
		return
	}
	defer func() {
		f.Close()
		if !x.TestingT.Failed() {
			os.Remove(f.Name())
		}
	}()

	diff, err := x.computeDiff(
		location(a.Output), bytes.NewReader(a.Output.Data),
		f.Name(), f,
	)
	if err != nil {
		x.TestingT.Log("unable to generate diff:", err)
		x.TestingT.Fail()
		return
	}

	if diff == nil {
		x.logSection("OUTPUT", a.Output.Data, location(a.Output))
	} else {
		x.logSection("OUTPUT DIFF", diff, "")
		x.TestingT.Fail()
	}
}
