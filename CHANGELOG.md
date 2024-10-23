# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

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

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
