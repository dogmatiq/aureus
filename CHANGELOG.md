# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [0.2.0] - 2024-10-20

### Added

- Add `Input` and `Output` interfaces to replace the `Content` and
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

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
