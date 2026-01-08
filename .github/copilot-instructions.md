# gohome AI Agent Instructions

## Project Overview

**gohome** is a Git activity aggregator CLI tool written in Go that scans workspace directories for git repositories and generates formatted daily standup reports from commit history. It follows Conventional Commits parsing and supports multiple output formats with clipboard integration.

## Architecture

### Core Pipeline (cmd/gohome/main.go)
The application follows a functional pipeline pattern:
1. **Config Load** → 2. **Scanner** → 3. **Git Client** → 4. **Parser** → 5. **Renderer**

```
main() → Load() → ScanGitRepos() → GetLogs() → Parse() → Print() → Clipboard
```

### Key Components

- **config** (`internal/config/`): Dual-source configuration (JSON file + CLI flags). File defaults at `~/.gohome.json` merged with flag overrides.
- **scanner** (`internal/scanner/`): Shallow directory scan (1-level deep) to discover `.git` folders. Skips `.git`, `.vscode`, `.idea`.
- **git** (`internal/git/`): Executes `git log` commands with sanitized inputs (regex-based injection prevention).
- **parser** (`internal/parser/`): Regex-based Conventional Commits parser extracting type/scope/message + emoji detection.
- **renderer** (`internal/renderer/`): Dual-format output (text/table) with preset styles (normal/markdown/nature/tech).
- **spinner** (`internal/spinner/`): Custom terminal spinner with configurable frames and intervals.

## Critical Patterns

### Configuration Precedence
CLI flags **override** JSON file values. Exception: time period flags (`--days`, `--hours`, etc.) use group mutual exclusion—if ANY time flag is set via CLI, ALL file time values are ignored.

Example from [config.go](internal/config/config.go#L236-L250):
```go
isTimeSetByUser := checkTimeFlags(userSetFlags)
if isTimeSetByUser {
    // Ignore ALL time values from file
} else {
    // Use ALL time values from file
}
```

### Security: Input Sanitization
All git command arguments are sanitized via regex before shell execution to prevent command injection:
```go
func sanitizeInput(input string) string {
    re := regexp.MustCompile(`[^a-zA-Z0-9\s._@-]+`)
    return re.ReplaceAllString(input, "")
}
```
See [git/client.go](internal/git/client.go#L25-L29)

### Entity Design
Two core entities in [entity/entity.go](internal/entity/entity.go):
- **Commit**: Parsed git log entry (Raw, Type, Scope, Message, Icon)
- **Task**: Manual/recurring task (supports `Enabled` flag for JSON persistence filtering)

Static tasks from `~/.gohome.json` are filtered by `Enabled: true`. CLI tasks (`-t`) are always shown.

## Development Workflows

### Build System (Makefile)
```bash
make build    # Compile to bin/gohome with version injection
make install  # Install to $GOPATH/bin
make test     # Run unit tests
make lint     # Run golangci-lint
```

Version info injected via LDFLAGS at build time:
```makefile
LDFLAGS=-ldflags "-X github.com/anIcedAntFA/gohome/internal/version.Version=$(VERSION) ..."
```

### Testing & Quality
- Use `go test -v ./...` for all packages
- Linting enforced via golangci-lint in CI (`.github/workflows/`)
- Security scanning: `#nosec` comments required for justified exclusions (e.g., validated file paths)

### Release Process
Automated via GoReleaser:
1. Tag version: `git tag v1.x.x`
2. Push: `git push origin v1.x.x`
3. GitHub Actions triggers cross-platform builds (Linux/Windows/macOS amd64/arm64)

See [RELEASE_GUIDE.md](RELEASE_GUIDE.md) and [.github/workflows/release.yml](.github/workflows/release.yml)

## Code Conventions

### Package Structure
- Use `internal/` for non-exported packages
- Single-responsibility services: each package exports one primary type (`Client`, `Service`, `Printer`, `Spinner`)
- Constructor pattern: `New<Type>()` functions (e.g., `git.NewClient()`)

### Error Handling
- Use `log.Fatal()` for unrecoverable errors in main flow
- Return `error` for recoverable operations (e.g., `scanner.ScanGitRepos()`)
- Silent failures for optional features (e.g., missing commits return empty slice)

### Spinner Usage Pattern
Always wrap long operations with spinners for UX:
```go
sp := spinner.New("🔍 Scanning repositories...").
    WithFrames(spinner.PacmanGhost).
    WithInterval(100 * time.Millisecond)
sp.Start()
// ... operation ...
sp.Stop()
```

See [spinner/spinner.go](internal/spinner/spinner.go) for custom frames.

### CLI Flag Design
- Support both long (`--flag`) and short (`-f`) forms
- Use `flag.Var()` for custom types (e.g., `StringSlice` for repeating `-t` flags)
- Implement `flag.Usage` override for formatted help via `tabwriter`

## External Dependencies

- **tablewriter** (`github.com/olekukonko/tablewriter`): Rich table formatting with preset styles
- **Standard library only** for core logic (no external frameworks)
- System clipboard via `internal/sys/clipboard.go` (platform-specific implementations)

## Current Limitations (from ROADMAP.md)

- No concurrency for multi-repo scanning yet (planned Fan-out/Fan-in pattern)
- No `--verbose` or `--quiet` flags
- No commit type filtering or directory exclusion patterns
- Limited unit test coverage

## Key Files to Reference

- [cmd/gohome/main.go](cmd/gohome/main.go): Main entry point and pipeline orchestration
- [internal/config/config.go](internal/config/config.go): Configuration loading logic and precedence rules
- [internal/parser/parser.go](internal/parser/parser.go): Conventional Commits regex and emoji extraction
- [Makefile](Makefile): Build commands and version injection
- [ROADMAP.md](ROADMAP.md): Feature status and planned enhancements

## When Making Changes

1. **Adding CLI flags**: Update both `Load()` function and `SaveToFile()` for JSON persistence
2. **Modifying git commands**: Ensure sanitization in `git.Client` methods
3. **New output formats**: Extend `renderer.Printer` with new `print*()` methods
4. **Version updates**: Use `make build` to inject version—never hardcode in source
5. **New spinner animations**: Add to `spinner/frames.go` FrameSet constants
