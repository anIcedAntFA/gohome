# Version Package

Provides build information and version formatting for the gohome CLI.

## Overview

This package handles version display logic with different formats for production releases vs development builds. It supports multiple build methods:

- **GoReleaser** (production releases) - ldflags injection
- **Local builds** (`make build`) - ldflags with VCS info
- **go install** - runtime VCS info from build metadata

## Version Format Logic

### Decision Flow

```
┌─────────────────────────────────────┐
│  Get version, commit, date          │
│  (from ldflags or VCS info)         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Is Semantic Version?               │
│  (has dots: v1.0.1, 2.3.0)          │
└──────────┬────────────┬─────────────┘
           │ YES        │ NO
           ▼            ▼
   ┌──────────┐   ┌─────────────────┐
   │ Clean    │   │ Has build info? │
   │ Format   │   │ (commit/date)   │
   └──────────┘   └────┬────────┬───┘
                       │ YES    │ NO
                       ▼        ▼
                  ┌────────┐  ┌────────┐
                  │ Full   │  │ Simple │
                  │ Details│  │ Format │
                  └────────┘  └────────┘
```

### Output Examples

| Build Method | Version Input | Output Format |
|--------------|---------------|---------------|
| **Production** (GoReleaser) | `v1.0.1` | `gohome v1.0.1` |
| **Dev Build** (`make build`) | `25bd8dd-dirty` | `gohome 25bd8dd-dirty (commit: 25bd8dd, built: 2026-01-10T15:10:26Z)` |
| **go install @latest** | `v1.0.0` | `gohome v1.0.0` |
| **go install @main** | `abc1234` | `gohome abc1234 (commit: abc1234, built: ...)` |

## Implementation Details

### Semantic Version Detection

The `isSemanticVersion()` function distinguishes between proper version tags and commit hashes:

**Semantic Versions (clean format):**
- Contains dots: `1.0.1`, `v2.3.0`, `v1.0.0-beta`
- Does NOT have dash at position 6 (commit hash pattern)

**Commit Hashes (full details):**
- 7-character hash + optional suffix: `abc1234-dirty`
- No dots, dash at position 6
- Literal "dev" string

### Build Info Sources (Priority Order)

1. **ldflags** (highest priority) - injected at build time:
   ```go
   -X github.com/anIcedAntFA/gohome/internal/version.Version={{.Version}}
   -X github.com/anIcedAntFA/gohome/internal/version.Commit={{.Commit}}
   -X github.com/anIcedAntFA/gohome/internal/version.Date={{.Date}}
   ```

2. **VCS info** (fallback) - from `runtime/debug.ReadBuildInfo()`:
   - `info.Main.Version` - for `go install github.com/user/repo@v1.0.0`
   - `vcs.revision` - git commit hash
   - `vcs.time` - commit timestamp

3. **Defaults** (last resort):
   - Version: `"dev"`
   - Commit: `"none"`
   - Date: `"unknown"`

## Functions

### Public API

- **`String() string`**  
  Returns full formatted version string with smart formatting based on build type.  
  Used by: `gohome --version`, `gohome -v`

- **`Short() string`**  
  Returns short version string (version only, no build details).  
  Reserved for future use (e.g., prompts, headers).

### Internal Helpers

- `getVersion() string` - retrieves version from ldflags or VCS
- `getCommit() string` - retrieves commit hash from ldflags or VCS
- `getDate() string` - retrieves build date from ldflags or VCS
- `isSemanticVersion(v string) bool` - checks if version is a proper semantic version tag

## Usage

```go
import "github.com/anIcedAntFA/gohome/internal/version"

// Full version string
fmt.Println(version.String())
// Output (production): "gohome v1.0.1"
// Output (dev): "gohome abc1234-dirty (commit: abc1234, built: 2026-01-10)"

// Short version string
fmt.Println(version.Short())
// Output: "gohome v1.0.1" or "gohome dev"
```

## Design Rationale

### Why Different Formats?

**Production releases** should be clean and professional:
- Users don't need to see internal commit hashes
- Shorter output for scripts and automation
- Matches industry standards (Docker, kubectl, etc.)

**Development builds** need full context:
- Developers need commit hash for debugging
- Build timestamp helps identify specific builds
- Clear indication that it's not a release version

### Why Check for Dots?

Semantic versions **always** have dots (`major.minor.patch`), while commit hashes never do. This is a reliable heuristic:

- ✅ `v1.0.1` → has dot → semantic version
- ✅ `2.3.0-beta` → has dot → semantic version
- ❌ `abc1234` → no dot → commit hash
- ❌ `25bd8dd-dirty` → no dot → dirty commit hash

### Edge Cases Handled

1. **Commit hash starting with digit**: `25bd8dd-dirty`
   - Checked dash position to distinguish from `2.5.0-beta`

2. **Version without 'v' prefix**: `1.0.1`
   - Strips prefix before checking, supports both formats

3. **go install without tag**: Uses short commit hash from VCS

4. **No build info available**: Falls back to "dev" gracefully

## Testing

Run unit tests:
```bash
go test -v ./internal/version/...
```

Test different scenarios:
```bash
# Dev build
make build && ./bin/gohome --version

# Simulated production
go build -ldflags "-X ...Version=v1.0.1" -o /tmp/test ./cmd/gohome
/tmp/test --version

# go install
go install github.com/anIcedAntFA/gohome/cmd/gohome@latest
gohome --version
```

## Future Enhancements

- [ ] Add JSON output format for scripting: `gohome version --json`
- [ ] Support prerelease tags: `v1.0.0-alpha.1`
- [ ] Add build metadata: `v1.0.0+build.123`
- [ ] Version comparison utilities (for update checks)
