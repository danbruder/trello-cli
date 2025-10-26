# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.2] - 2025-10-25

## [1.0.1] - 2025-10-25

### Added
- Comprehensive unit tests for all major modules:
  - Config utility functions (maskString, getConfigPath)
  - Batch operation validation for all 7 operation types
  - Formatter methods for labels, checklists, members, and attachments
  - Context optimizer edge cases and boundary conditions
  - Client initialization tests

### Fixed
- Fixed `TruncateText` function to handle negative maxLen values (prevented panic)
- Updated tests to gracefully handle optional config files

### Improved
- Test coverage now includes 1,500+ lines of comprehensive tests
- Added edge case testing for token limiting, field filtering, and text truncation
- All tests passing across all modules
## [1.0.0] - 2025-10-25

### Added
- Initial release with comprehensive Trello API support
- LLM-optimized output formats (Markdown and JSON)
- Batch operations support
- Context optimization with token limits and field filtering
- Cross-platform binary builds (Linux, macOS, Windows)
- Docker image distribution
- Homebrew tap distribution for macOS (`brew tap danbruder/tap && brew install trello-cli`)
- GitHub Actions for automated releases

### Features
- Full CRUD operations on boards, lists, cards, labels, checklists, members, and attachments
- Flexible authentication (environment variables, config file, command-line flags)
- Scripting support with quiet mode
- Comprehensive error handling with appropriate exit codes
- Configuration management with config command

[Unreleased]: https://github.com/danbruder/trello-cli/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/danbruder/trello-cli/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/danbruder/trello-cli/releases/tag/v1.0.0
