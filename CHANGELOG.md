# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [0.2.7] - 2024-11-26

### Changed

- File-based tests may now omit the leading "atom" in filenames before the
  `.input` or `.output` part.

## [0.2.6] - 2024-11-25

### Added

- Added a suggested `go test` command for re-running a failed test without
  blessing it. There is nothing specific to Aureus about this command, but it's
  useful considering the `-run` flag's pattern syntax can make it tricky to
  isolate a single test.

### Changed

- `-aureus.bless` now prevents the blessed test from being marked as a failure.
  This means you don't see the output of the test, so the `go test` command that
  is suggested when a test may be blessed now includes `-v`.
- Suggested `go test` commands now include `-count 1` to prevent the test cache
  from preventing a re-run of the test, even if the input or output files have
  changed.
- A passing test now renders the golden file in dim yellow ("gold") color.
  Failed tests continue to use a colorized diff.
- Changed how suggested commands are rendered to make them occupy less visual
  space.

### Fixed

- Fixed issue that could cause loss of trailing characters of fenced code blocks
  when blessing markdown tests.

## [0.2.5] - 2024-11-19

### Fixed

- Fixed order of parameters when suggesting `go test` commands to run.
- Attempt to suggest a single package to run with `-aureus.bless`, instead of
  `./...`. The latter works, but can produce error output messages if some
  packages don't use Aureus.
- Fixed escaping of shell parameters.

## [0.2.4] - 2024-11-18

### Added

- Added the `-aureus.bless` test flag, which "blesses" the output of a failing
  test, accepting it as correct for future test runs. For file-based tests it
  replaces the entire output file with the failing output. For Markdown tests it
  replaces the code within the fenced code block.

## [0.2.3] - 2024-11-13

### Fixed

- Removed duplicate "unable to generate output" error text.

## [0.2.2] - 2024-10-23

### Changed

- `Run()` no longer separates the tests built from flat-files from those built
  from Markdown documents. This provides a less surprising test heirarchy when
  both types of tests are present in the same directory, and eliminates an empty
  test when only one type of test is present.
- Tests built from Markdown documents that contain a single top-level heading at
  the top of the document are now named after that heading, rather than the
  filename.
- Matrix tests with inputs and outputs from codeblocks under Markdown headings
  are now named after those headings, rather than the filename/line number.
- Matrix tests that only have a single input or a single output no longer
  include the name of that single input/output, producing shorter test names.
- Tests with inputs and outputs sourced from the same file now only include the
  filename in the test name once.
- Diff output is now colorized.

## [0.2.1] - 2024-10-20

### Fixed

- Removed empty "extra" section in `BEGIN OUTPUT DIFF` log message.
- Added missing space to `END <section>` log messages.

## [0.2.0] - 2024-10-20

### Added

- Added `Input` and `Output` interfaces to replace the `Content` and
  `ContentMetaData` types previously passed to output generators.
- Generated output is now written to a temporary file when a test fails, to
  make it easier to debug.

### Changed

- **[BC]** Changed the signature of the output generation passed to `Run()`. It
  now accepts the `*testing.T` value (or equivalent) and the new `Input` and
  `Output` interfaces.
- Improved test output.
- Improved diff formatting.

## [0.1.0] - 2024-02-19

- Initial release

<!-- references -->

[Unreleased]: https://github.com/dogmatiq/aureus
[0.1.0]: https://github.com/dogmatiq/aureus/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.0
[0.2.1]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.1
[0.2.2]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.2
[0.2.3]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.3
[0.2.4]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.4
[0.2.5]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.5
[0.2.6]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.6
[0.2.7]: https://github.com/dogmatiq/aureus/releases/tag/v0.2.7

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
