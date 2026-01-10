# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.2] - 2026-01-10

### Added

- Multi-platform package manager support in install script (apt, dnf, yum, zypper, apk, pacman, brew)
- Download verification with file type checking
- Comprehensive shell configuration examples (bash/zsh/fish) in README
- Internal documentation for version package with design rationale

### Fixed

- Install script non-interactive mode (curl | bash) now auto-accepts upgrade prompts
- Architecture parsing in install script (correctly extracts `x86_64` instead of `64`)
- PATH priority conflicts by auto-removing dev builds from `$GOPATH/bin`
- Config flag syntax error preventing config file loading
- golangci-lint issues:
  - emptyStringTest: use `test == ""` instead of `len(test) == 0`
  - gocyclo: reduced cyclomatic complexity of `Load()` from 17 to 5 via helper extraction

### Changed

- **Version format differentiation:**
  - Production releases: clean format `gohome v1.0.2` (no build details)
  - Development builds: full format `gohome abc1234 (commit: abc1234, built: 2026-01-10)`
- Refactored `config.Load()` into smaller helper functions for better maintainability
- Install script now warns and verifies correct binary location in PATH

### Documentation

- Enhanced README installation section with:
  - Install script features and behavior
  - PATH configuration best practices for all major shells
  - Version format explanation (production vs dev)
  - Collapsible shell config sections
- Added comprehensive `internal/version/README.md` with:
  - Decision flow diagram
  - Build method comparison table
  - Semantic version detection logic
  - Testing guide and examples

## [1.0.1] - 2026-01-08

### Fixed

- Version command now correctly displays version information
- Fixed CI/CD workflows to inject version at build time
- Improved version display for `go install` users (cleaner output)
- Fixed prealloc lint warnings with proper slice preallocation
- Fixed Windows build compatibility in CI workflows

### Changed

- Refactored version handling into dedicated `internal/version` package
- Enhanced version detection with VCS fallback for go install users

### Documentation

- Added comprehensive VERSIONING.md guide
- Updated README with version flag usage

## [1.0.0] - 2026-01-08

### Added

- Release automation with GoReleaser
- Version support with `--version` / `-v` flag
- Universal installation script (curl|sh)
- GitHub Actions workflow for automated releases
- Multi-platform builds (Linux, macOS, Windows)
- Comprehensive release documentation (RELEASE_GUIDE, RELEASE_CHECKLIST, SUMMARY)

### Changed

- Improved flag parsing to support version checking

### Fixed

- Flag parsing conflict between version flag and config flags

## [0.1.0] - 2026-01-07

### Added

- Git commit aggregation and reporting
- Custom tasks support (static from config + dynamic from CLI)
- Multiple output formats (text, table, markdown)
- Copy to clipboard functionality
- Loading spinner for better UX
- Config file persistence (~/.gohome.json)
- Flexible time period options (hours, days, weeks, months, years)
- Icon and scope display options
- Multiple repository support
- Conventional commits parsing

### Documentation

- README with usage examples
- ROADMAP with development milestones
- Release guides and checklists

[Unreleased]: https://github.com/anIcedAntFA/gohome/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/anIcedAntFA/gohome/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/anIcedAntFA/gohome/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/anIcedAntFA/gohome/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/anIcedAntFA/gohome/releases/tag/v0.1.0
