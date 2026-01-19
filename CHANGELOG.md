# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0-beta.2] - 2026-01-19

### Fixed

- **Config Command Error Handling**
  - Add consistent error formatting with ❌ prefix across all commands
  - `gohome config <unknown>` now shows proper error instead of just help
  - Add `SilenceErrors`/`SilenceUsage` to both config and report commands
  - Add Args validation to config command for unknown subcommands

- **Config List Output Enhancement**
  - Add ⚙️ emoji to 'Current Configuration' header
  - Add 📄 emoji to config file path display
  - Remove unused table options (nature and tech type)

- **NPM Prerelease Publishing**
  - Fix NPM publish workflow for beta/alpha/rc versions
  - Auto-detect prerelease versions and apply correct `--tag` flag
  - Prevent beta versions from overwriting `latest` NPM tag
  - Support conditional publishing: prereleases use custom tags, stable uses `latest`

### Documentation

- Add comprehensive release checklist template (`.github/RELEASE_CHECKLIST_CURRENT.md`)
- Add release notes template for v1.3.0-beta.1
- Add NPM prerelease fix documentation (`.github/NPM_PRERELEASE_FIX.md`)
- Update README with improved examples
- Update CLI_GUIDE with latest command usage
- Add beta release testing script (`scripts/test-beta-release.sh`)
- Update quickstart demo with better visuals

## [1.3.0-beta.1] - 2026-01-18

### Added

- **Cobra/Viper Framework Integration**
  - Modern CLI with subcommands: `report`, `config`, `version`, `completion`
  - Improved flag binding and validation
  - Better organized help text with ASCII art logo
  
- **Multi-Format Configuration Support**
  - Support for JSON, YAML, and TOML config files
  - Config file location: `~/.config/gohome/config.{json,yaml,toml}`
  - Backward compatible with `~/.gohome.json`
  
- **Environment Variables**
  - Configure via `GOHOME_*` environment variables
  - All CLI flags can be set through env vars
  - Priority: CLI flags > env vars > config file > defaults
  
- **Config Management Subcommand**
  - `gohome config init` - Initialize config with defaults
  - `gohome config show` - Display current configuration
  - `gohome config edit` - Edit config in default editor
  
- **Shell Completions**
  - Auto-completion for bash, zsh, fish, powershell
  - `gohome completion <shell>` to generate completion scripts
  - Improves CLI discoverability and productivity
  
- **Comprehensive Test Suite**
  - 100% coverage for critical packages (parser, git, renderer)
  - Security-focused tests (command injection prevention)
  - Overall project coverage: 45.1% (up from 14.6%)

### Changed

- **Configuration File Location**
  - New default: `~/.config/gohome/config.{json,yaml,toml}`
  - Old location (`~/.gohome.json`) still supported for backward compatibility
  
- **CLI Architecture**
  - Refactored to use Cobra/Viper for better maintainability
  - Subcommand structure for clearer UX
  - Consistent error handling across all commands

### Fixed

- All 47 linting errors resolved (gocritic, godot, gosec, misspell, revive)
- Improved input sanitization for security
- Better error messages with context

### Documentation

- Added comprehensive migration guide (v1.2 → v1.3)
- Updated README with v1.3 features
- Enhanced CONTRIBUTING with testing guidelines
- Updated CLI_GUIDE with new commands

## [1.2.0] - 2026-01-15

### Added

- **Recursive scanning with configurable depth** (`--max-depth` flag)
  - Scanner now supports multi-level directory structures (default 2 levels deep)
  - Properly detects repos in nested paths like `github.com/{org}/{repo}`
  - Configurable via `--max-depth <int>` flag or `max_depth` in JSON config
  - Automatically defaults to 2 levels if maxDepth <= 0
  - Implements `scanRecursive()` helper function with depth tracking
  - Fixes issue where only shallow (1-level) scanning was performed
  - See PR #11 for implementation details

- **Git LFS (Large File Storage) for demo media files**
  - Migrated demo GIF files from regular Git to LFS storage
  - Files tracked: `docs/demos/*.{gif,png,mp4,webm}`
  - Reduced Git history size: 4.3MB → 264 bytes (LFS pointers)
  - Scalable solution for adding unlimited demo content
  - Proper `.gitattributes` configuration following VHS project best practices
  - CI/CD workflows updated with `lfs: true` to fetch files during builds

- **Comprehensive project documentation**
  - `CONTRIBUTING.md`: Development setup, workflow, coding standards, commit conventions
  - `SECURITY.md`: Vulnerability reporting, security policies, best practices
  - `docs/GIT_LFS_GUIDE.md`: Complete Git LFS setup and usage guide (320+ lines)
  - `docs/RELEASE_CHECKLIST.md` and `docs/RELEASE_GUIDE.md` moved to `docs/` folder for better organization

- **WSL2 clipboard support**
  - Fixed clipboard functionality for Windows Subsystem for Linux 2
  - Refactored clipboard detection to use switch statement for better maintainability
  - Properly handles WSL2 environment detection and `clip.exe` integration

- **Comprehensive scanner unit tests**
  - Added 4 test suites: `TestIsGitRepo`, `TestShouldSkip`, `TestScanRecursive`, `TestScanGitRepos`
  - Covers depth scenarios (0, 1, 2, 3 levels), special directories, edge cases
  - Table-driven tests with comprehensive coverage
  - Test helper `createGitRepo()` for consistent test setup

### Changed

- **Documentation reorganization**
  - Moved `RELEASE_CHECKLIST.md` and `RELEASE_GUIDE.md` from root to `docs/` folder
  - Updated all cross-references in `ROADMAP.md` and `.github/copilot-instructions.md`
  - Improved project structure for better documentation discoverability

- **CI/CD improvements**
  - Release workflow now properly fetches LFS files during checkout
  - Build workflow updated with LFS support
  - Automated npm version synchronization during release publishing

### Fixed

- **Scanner not detecting nested repositories**
  - Fixed shallow scanning limitation (was only 1 level deep)
  - Now properly discovers repos in `github.com/{org}/{repo}` structures
  - Respects configurable depth setting with proper validation

- **Git LFS workflow errors**
  - Resolved "dirty state" detection during GoReleaser builds
  - Fixed LFS pointer files being checked out instead of actual media files
  - Proper workflow configuration prevents LFS-related release failures

### Documentation

- Updated `.github/copilot-instructions.md` with:
  - Recursive scanner implementation details
  - Git LFS configuration and usage patterns
  - New documentation file references
  - Enhanced "When Making Changes" guidelines
- Updated `README.md` with Git LFS notes (if applicable)
- Updated `ROADMAP.md` with completed features and future distribution plans

## [1.1.0] - 2026-01-13

### Added

- **Multi-branch support** (`--all-branches` / `-b` flag)
  - Include commits from all local branches using `git log --branches`
  - Solves the problem where commits on unmerged feature branches are invisible
  - No need to pull from remote or wait for PR merges
  - Git automatically handles commit deduplication
  - Configurable via CLI flag or JSON config file (`all_branches: true`)
  - Perfect for generating standup reports before merging PRs

- **Branch filtering** (`--branch <name>` flag)
  - Filter commits from a specific branch without checking it out
  - Useful for reviewing work on feature branches
  - Mutually exclusive with `--all-branches` (branch filter takes precedence)
  - Branch names are sanitized to prevent command injection
  - Configurable via CLI flag or JSON config file (`branch: "branch-name"`)

### Changed

- Refactored `mergeConfigs` function to reduce cyclomatic complexity (17 → under 15)
  - Extracted `mergeTimePeriods`, `mergeStringFlags`, `mergeBooleanFlags`, `mergeTasks` helper functions
  - Improved code maintainability and readability
  - No behavioral changes, pure refactoring

### Documentation

- Updated README.md with comprehensive branch features documentation
  - Added `--all-branches` and `--branch` to Flags Reference table
  - Added usage examples (6️⃣ Include All Local Branches, 7️⃣ Filter by Specific Branch)
  - Added important note about default behavior (shows current branch only)
  - Updated Features section to mention branch support
- Updated ROADMAP.md
  - Moved multi-branch and branch filtering features to Phase 1 (completed)
  - Marked features as implemented in filtering section
- Updated example config with `all_branches` and `branch` fields

## [1.0.4] - 2026-01-12

### Added

- **AUR (Arch User Repository) support**
  - Automated publishing to AUR via GoReleaser
  - Build from source following Arch packaging guidelines
  - Installation: `yay -S gohome` or `paru -S gohome`
  - Includes proper PKGBUILD with build flags and version injection
  - See [docs/AUR_SETUP.md](docs/AUR_SETUP.md) for setup guide

### Fixed

- **NPM package installation issues**
  - Fixed `Cannot find module './package.json'` error in postinstall script
  - Fixed incorrect `binDir` path (was creating `scripts/bin` instead of root `bin/`)
  - Fixed architecture mapping: `amd64` → `x86_64` to match GoReleaser output
  - Fixed file extension: removed Windows `.zip` handling, all platforms use `.tar.gz`
  - Simplified extraction logic: unified `tar` command for all platforms (Windows 10+ supported)
  - Added prerelease version handling: strips `-beta`, `-hotfix` suffixes for GitHub release URL

### Changed

- **NPM publishing automation**
  - Migrated from token-based to OIDC trusted publishing for enhanced security
  - Removed long-lived `NPM_TOKEN` requirement
  - Added automatic provenance attestation generation
  - Updated npm CLI to latest version (>= 11.5.1) in CI workflow

## [1.0.3] - 2026-01-11

### Added

- **NPM distribution support** via GoReleaser
  - Package published as `@ngockhoi96/gohome` on npm registry
  - Installation: `npm install -g @ngockhoi96/gohome`
  - Supports npx usage: `npx @ngockhoi96/gohome`
  - Automated publishing via GitHub Actions with OIDC trusted publishing
  - Structured package with bin/ and scripts/ folders
  - Wrapper script for seamless binary execution
  - Test suite for package validation
- Windows PowerShell installation script (`scripts/install.ps1`)
  - Auto-detect architecture (x64/arm64/x86)
  - Install to `%LOCALAPPDATA%\Programs\gohome`
  - Automatically add to user PATH
  - Clean up conflicting GOPATH binaries

### Changed

- Reorganized installation scripts into `scripts/` folder
- Updated documentation with PowerShell installation examples
- Enhanced shell configuration guide with PowerShell PATH management
- Updated release notes template to include npm installation method

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
