<div align="center">

# Aureus

A test runner for executing "golden file" tests in Go.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/aureus)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/aureus.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/aureus/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/aureus/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/aureus/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/aureus/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/aureus)

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

[`testdata`]: testdata
[`run_test.go`]: run_test.go
[readme source]: https://github.com/dogmatiq/aureus/blob/main/README.md?plain=1
[fenced code blocks]: https://spec.commonmark.org/0.31.2/#fenced-code-blocks
