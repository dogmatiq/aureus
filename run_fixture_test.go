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
		document := `~~~sql au:fixture au:group=example role=first
first
~~~

~~~text au:input au:group=example
VALUE
~~~

~~~text au:input au:group=example
value
~~~

~~~text au:output au:group=example
first|second|value
~~~

~~~text au:output au:group=example
first|second|value
~~~

~~~text au:fixture au:group=example role=second
second
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
				if first.Language() != "sql" || first.Attributes()["role"] != "first" {
					t.Fatal("expected the first fixture's language and attributes")
				}
				second := a.Fixtures[1]
				if second.Language() != "text" || second.Attributes()["role"] != "second" {
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
					strings.ToLower(strings.TrimSpace(string(input))),
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
first
~~~

~~~text au:fixture au:group=second
second
~~~

~~~text au:input au:group=first
value
~~~

~~~text au:output au:group=first
first|value
~~~

~~~text au:input au:group=second
value
~~~

~~~text au:output au:group=second
second|value
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
unavailable fixture
~~~

~~~text au:input au:group=example
input
~~~

~~~text au:output au:group=example
output
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
first fixture
~~~

~~~text au:input au:group=example
input
~~~

~~~text au:output au:group=example
old output
~~~

~~~text au:fixture au:group=example
second fixture
~~~
`
		dir := writeTestDocument(t, document)

		Run(
			t,
			func(_ *testing.T, a Assertion) error {
				_, err := io.WriteString(a.Output, "new, longer output\n")
				return err
			},
			FromDir(dir),
			Bless(true),
		)

		got, err := os.ReadFile(filepath.Join(dir, "example.md"))
		if err != nil {
			t.Fatal(err)
		}

		want := strings.Replace(document, "old output\n", "new, longer output\n", 1)
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
