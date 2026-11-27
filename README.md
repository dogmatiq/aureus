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
respectively. Files named `<group>.fixture[.<extension>]` provide fixtures for
each test in the same group.

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
group the inputs and outputs. Blocks annotated with `au:fixture` supply content
used to prepare each test in the named group.

This file is itself an example of a Markdown-based test. It confirms the
behavior of a basic JSON pretty-printer. This fixture configures the formatter
to use four spaces of indentation:

```json au:fixture au:group=json-pretty-printer role=format
{ "indent": "    " }
```

Given this unformatted JSON value:

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

[`testdata`]: testdata
[`run_test.go`]: run_test.go
[readme source]: https://github.com/dogmatiq/aureus/blob/main/README.md?plain=1
[fenced code blocks]: https://spec.commonmark.org/0.31.2/#fenced-code-blocks
