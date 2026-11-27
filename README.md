<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://dogmatiq.io/res/img/dynamic/repo/aureus/banner/dark.svg">
  <img alt="Dogmatiq aureus" src="https://dogmatiq.io/res/img/dynamic/repo/aureus/banner/light.svg" height="44">
</picture>

A test runner for executing "golden file" tests in Go.

[![API reference](https://dogmatiq.io/res/img/dynamic/repo/aureus/badge/docs.svg)](https://pkg.go.dev/github.com/dogmatiq/aureus)
[![Latest version](https://dogmatiq.io/res/img/dynamic/repo/aureus/badge/semver.svg)](https://github.com/dogmatiq/aureus/releases)
[![Build status](https://dogmatiq.io/res/img/dynamic/repo/aureus/badge/build.svg)](https://github.com/dogmatiq/aureus/actions/workflows/ci.yml)
[![Code coverage](https://dogmatiq.io/res/img/dynamic/repo/aureus/badge/coverage.svg)](https://codecov.io/github/dogmatiq/aureus)

</div>

## What is a golden file?

A "golden file" is a file that represents the expected output of some test. When
the test is executed, the output of the system under test is compared to the
content of the golden file — if they differ, the test fails.

## What does Aureus do?

Aureus recursively scans directories for golden file tests expressed either as
flat files, or as code blocks within Markdown documents. By default it scans the
`testdata` directory within the current working directory.

### Flat files

Files with names matching the `<group>.input[.<extension>]` and
`<group>.output[.<extension> ` patterns are treated as test inputs and outputs,
respectively.

For each pair of input and output files with the same `<group>` prefix, a
user-defined function is invoked. The function must produce output that matches
the content of the output file, otherwise the test fails.

The [`run_test.go`] illustrates how to use Aureus to execute the flat-file tests
in the [`testdata`] directory.

### Markdown documents

As an alternative to (or in combination with) flat-file tests, Aureus can load
inputs and outputs from [fenced code blocks] within Markdown documents.

Code blocks annotated with the `au:input` or `au:output` attribute are treated
as a test input or output, respectively. The `au:group` attribute is used to
group the inputs and outputs.

This file is itself an example of a Markdown-based test. It confirms the
behavior of a basic JSON pretty-printer. Given this unformatted JSON value:

```json au:input au:group=json-pretty-printer
{ "one": 1, "two": 2 }
```

We expect our formatter function to produce the following output:

```json au:output au:group=json-pretty-printer
{
  "one": 1,
  "two": 2
}
```

View the [README source] to see how the code blocks are annotated for use with
Aureus, and [`run_test.go`] to see how to execute the
tests.

### Fixtures

Fixtures supply content used to prepare each assertion. They belong to a named
group, alongside its inputs and expected outputs.

For example, this fixture configures the JSON pretty-printer's indentation:

```json au:fixture au:group=custom-indent role=format
{"indent": "    "}
```

Given this input:

```json au:input au:group=custom-indent
{"one": 1, "two": 2}
```

We expect four spaces of indentation:

```json au:output au:group=custom-indent
{
    "one": 1,
    "two": 2
}
```

The generator receives one `Assertion` containing `Fixtures`, `Input` and
`Output`. Each fixture exposes its content as a reader, plus its `Language()`
and custom `Attributes()`. Here, `role=format` is a custom attribute interpreted
by the pretty-printer. Aureus supplies the content; the generator decides how
to prepare the assertion and clean up any resources.

Every input/output combination in the group receives all its fixtures, with
fresh readers for each assertion. Fixtures do not add combinations to the test
matrix. Markdown fixtures are supplied in document order, and apply to their
entire group regardless of where they appear in the document.

Flat-file tests use `<group>.fixture[.<extension>]`. Custom attributes use the
same filename syntax as inputs and outputs, for example
`example.fixture.@role=format.json`. Multiple fixture files are supplied in
filename order. The group prefix is required.

Groups are local to a Markdown document or a directory of flat-file tests.
Fixtures are not inherited by other documents or directories. A fixture without
an assertion is an error. Skipping a fixture with `au:skip`, or a leading
underscore on its filename, skips every assertion that requires it.

### Updating existing generators

Generators now accept `func(t T, a Assertion) error` instead of separate input
and output arguments. Read from `a.Input` and write to `a.Output`. Existing
tests without fixtures receive an empty `a.Fixtures` slice.

[`testdata`]: testdata
[`run_test.go`]: run_test.go
[readme source]: https://github.com/dogmatiq/aureus/blob/main/README.md?plain=1
[fenced code blocks]: https://spec.commonmark.org/0.31.2/#fenced-code-blocks
