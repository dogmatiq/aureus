package aureus_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/dogmatiq/aureus"
)

func TestRun_fixtures(t *testing.T) {
	t.Run("it supplies fresh fixtures to every input and output combination", func(t *testing.T) {
		document := `~~~text au:fixture au:group=example foo=bar
FIXTURE 1
~~~

~~~text au:input au:group=example
INPUT
~~~

~~~text au:input au:group=example
input
~~~

~~~text au:output au:group=example
FIXTURE 1|FIXTURE 2|INPUT
~~~

~~~text au:output au:group=example
FIXTURE 1|FIXTURE 2|INPUT
~~~

~~~text au:fixture au:group=example foo=baz
FIXTURE 2
~~~
`
		dir := writeTestDocument(t, document)

		var calls int
		Run(
			t,
			func(t *testing.T, a Assertion) error {
				calls++

				if len(a.Fixtures) != 2 {
					t.Fatalf("expected two fixtures, got %d", len(a.Fixtures))
				}

				first := a.Fixtures[0]
				if first.Language() != "text" || first.Attributes()["foo"] != "bar" {
					t.Fatal("expected the first fixture's language and attributes")
				}
				second := a.Fixtures[1]
				if second.Language() != "text" || second.Attributes()["foo"] != "baz" {
					t.Fatal("expected the second fixture's language and attributes")
				}

				prefix, err := io.ReadAll(first)
				if err != nil {
					return err
				}
				suffix, err := io.ReadAll(second)
				if err != nil {
					return err
				}
				input, err := io.ReadAll(a.Input)
				if err != nil {
					return err
				}

				_, err = fmt.Fprintf(
					a.Output,
					"%s|%s|%s\n",
					strings.TrimSpace(string(prefix)),
					strings.TrimSpace(string(suffix)),
					strings.ToUpper(strings.TrimSpace(string(input))),
				)
				return err
			},
			FromDir(dir),
		)

		if calls != 4 {
			t.Fatalf("expected four assertion pairs, got %d", calls)
		}
	})
}

func TestRun_fixtureGroups(t *testing.T) {
	t.Run("it only supplies fixtures from the assertion's group", func(t *testing.T) {
		dir := writeTestDocument(t, `~~~text au:fixture au:group=first
FIXTURE 1
~~~

~~~text au:fixture au:group=second
FIXTURE 2
~~~

~~~text au:input au:group=first
INPUT
~~~

~~~text au:output au:group=first
FIXTURE 1|INPUT
~~~

~~~text au:input au:group=second
INPUT
~~~

~~~text au:output au:group=second
FIXTURE 2|INPUT
~~~
`)

		var calls int
		Run(
			t,
			func(t *testing.T, a Assertion) error {
				calls++
				if len(a.Fixtures) != 1 {
					t.Fatalf("expected one fixture from this group, got %d", len(a.Fixtures))
				}

				fixture, err := io.ReadAll(a.Fixtures[0])
				if err != nil {
					return err
				}
				input, err := io.ReadAll(a.Input)
				if err != nil {
					return err
				}

				_, err = fmt.Fprintf(
					a.Output,
					"%s|%s\n",
					strings.TrimSpace(string(fixture)),
					strings.TrimSpace(string(input)),
				)
				return err
			},
			FromDir(dir),
		)

		if calls != 2 {
			t.Fatalf("expected one assertion per group, got %d", calls)
		}
	})
}

func TestRun_skippedFixture(t *testing.T) {
	t.Run("it does not invoke the generator when a required fixture is skipped", func(t *testing.T) {
		dir := writeTestDocument(t, `~~~text au:fixture au:group=example au:skip
FIXTURE
~~~

~~~text au:input au:group=example
INPUT
~~~

~~~text au:output au:group=example
OUTPUT
~~~
`)

		Run(
			t,
			func(t *testing.T, _ Assertion) error {
				t.Fatal("the generator must not run without its required fixture")
				return nil
			},
			FromDir(dir),
		)
	})
}

func TestRun_blessWithFixtures(t *testing.T) {
	t.Run("it replaces only the expected output when fixtures surround an assertion", func(t *testing.T) {
		document := `~~~text au:fixture au:group=example
FIXTURE 1
~~~

~~~text au:input au:group=example
INPUT
~~~

~~~text au:output au:group=example
OUTPUT
~~~

~~~text au:fixture au:group=example
FIXTURE 2
~~~
`
		dir := writeTestDocument(t, document)

		Run(
			t,
			func(_ *testing.T, a Assertion) error {
				_, err := io.WriteString(a.Output, "LONGER OUTPUT\n")
				return err
			},
			FromDir(dir),
			Bless(true),
		)

		got, err := os.ReadFile(filepath.Join(dir, "example.md"))
		if err != nil {
			t.Fatal(err)
		}

		want := strings.Replace(document, "OUTPUT\n", "LONGER OUTPUT\n", 1)
		if string(got) != want {
			t.Fatalf("expected only the output block to change:\n%s", got)
		}
	})
}

func writeTestDocument(t testing.TB, document string) string {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "example.md"), []byte(document), 0600); err != nil {
		t.Fatal(err)
	}

	return dir
}
