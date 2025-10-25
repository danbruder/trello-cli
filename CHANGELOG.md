# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2025-10-25

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

## [1.0.0] - TBD

### Added
- Initial release

[Unreleased]: https://github.com/danbruder/trello-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/danbruder/trello-cli/releases/tag/v1.0.0
