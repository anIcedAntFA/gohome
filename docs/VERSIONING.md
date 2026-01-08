# Version Management Guide

## 📖 Table of Contents

- [Overview](#overview)
- [How Versioning Works](#how-versioning-works)
- [Version Display Scenarios](#version-display-scenarios)
- [Release Process](#release-process)
- [For Contributors](#for-contributors)
- [Troubleshooting](#troubleshooting)

---

## Overview

`gohome` uses **semantic versioning** (SemVer) with the format `vMAJOR.MINOR.PATCH`:

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

Example: `v1.2.3`

---

## How Versioning Works

### 🏗️ Build Time Version Injection

Version information is injected at **build time** using Go's `-ldflags`:

```bash
go build -ldflags "
  -X github.com/anIcedAntFA/gohome/internal/version.Version=v1.0.0
  -X github.com/anIcedAntFA/gohome/internal/version.Commit=abc123
  -X github.com/anIcedAntFA/gohome/internal/version.Date=2026-01-08T12:00:00Z
" ./cmd/gohome
```

### 📦 Version Sources (Priority Order)

1. **ldflags injection** (highest priority)
   - Used by: Local builds with `make build`, GoReleaser
   - Source: Git tags via `git describe --tags`

2. **VCS (Version Control System) info** (fallback)
   - Used by: `go install` without ldflags
   - Source: Embedded by Go toolchain from git metadata
   - Format: `v0.0.0-20260108065610-abc123def` (pseudo-version)

3. **Default values** (last resort)
   - Used when: No git info available
   - Values: `Version="dev"`, `Commit="none"`, `Date="unknown"`

### 🔍 Code Implementation

```go
// internal/version/version.go
func getVersion() string {
    // 1. Check ldflags (set during make build or GoReleaser)
    if Version != "dev" {
        return Version  // e.g., "v1.0.0"
    }

    // 2. Fallback to VCS info (for go install users)
    if info, ok := debug.ReadBuildInfo(); ok {
        if info.Main.Version != "" && info.Main.Version != "(devel)" {
            return info.Main.Version  // e.g., "v0.0.0-20260108..."
        }
        
        // 3. Try commit hash from VCS
        for _, setting := range info.Settings {
            if setting.Key == "vcs.revision" {
                return setting.Value[:7]  // Short hash
            }
        }
    }

    // 4. Last resort
    return "dev"
}
```

---

## Version Display Scenarios

### Scenario 1: Local Development (make build)

```bash
$ make build
$ ./bin/gohome --version
gohome 7a88856-dirty (commit: 7a88856, built: 2026-01-08T07:18:19Z)
```

**Explanation:**
- `7a88856` = current commit hash (no tag found)
- `-dirty` = uncommitted changes exist
- Makefile uses `git describe --tags --always --dirty`

### Scenario 2: Release Build (GoReleaser)

```bash
# After: git tag v1.0.0 && git push origin v1.0.0
# GitHub Actions runs GoReleaser

$ gohome --version  # Downloaded from GitHub Releases
gohome v1.0.0 (commit: abc123, built: 2026-01-08T12:00:00Z)
```

**Explanation:**
- GoReleaser detects tag `v1.0.0`
- Injects version via ldflags in `.goreleaser.yml`

### Scenario 3: Go Install (without tag)

```bash
$ go install github.com/anIcedAntFA/gohome/cmd/gohome@latest
$ gohome --version
gohome v0.0.0-20260108065610-7a88856 (commit: none, built: unknown)
```

**Explanation:**
- No ldflags injection (go install doesn't use them)
- Falls back to VCS pseudo-version
- Format: `v0.0.0-YYYYMMDDHHMMSS-commithash`

### Scenario 4: Go Install (with tag)

```bash
$ go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.0.0
$ gohome --version
gohome v1.0.0 (commit: none, built: unknown)
```

**Explanation:**
- Go modules resolves to tagged version
- VCS info contains `v1.0.0`
- Commit/date not available (not injected via ldflags)

---

## Release Process

### Step-by-Step Release

#### 1. Prepare Release Branch

```bash
# Update from main
git checkout main
git pull origin main

# Create release branch (optional, for hotfixes)
git checkout -b release/v1.1.0
```

#### 2. Update Version References

Update documentation if needed:
- `README.md` examples
- `CHANGELOG.md` entry

Commit changes:
```bash
git add .
git commit -m "📝 docs: prepare v1.1.0 release"
git push origin release/v1.1.0
```

#### 3. Create and Push Tag

```bash
# Create annotated tag (recommended)
git tag -a v1.1.0 -m "Release v1.1.0

Features:
- Add new feature X
- Improve performance of Y

Bug Fixes:
- Fix issue #42
"

# Push tag to trigger release workflow
git push origin v1.1.0
```

#### 4. GitHub Actions Automation

Once tag is pushed, GitHub Actions automatically:

1. **Trigger** `.github/workflows/release.yml`
2. **Run tests** with coverage
3. **Build** binaries for all platforms (via GoReleaser)
4. **Create** GitHub Release with:
   - Changelog (auto-generated)
   - Binaries (Linux, macOS, Windows × amd64, arm64)
   - Checksums
5. **Upload** coverage to Codecov

#### 5. Verify Release

Check GitHub Releases page:
```
https://github.com/anIcedAntFA/gohome/releases/tag/v1.1.0
```

Test installation:
```bash
# Method 1: Download binary
curl -sSL https://github.com/anIcedAntFA/gohome/releases/download/v1.1.0/gohome_1.1.0_Linux_x86_64.tar.gz | tar xz
./gohome --version

# Method 2: Go install
go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.1.0
gohome --version
```

#### 6. Merge to Main

```bash
git checkout main
git merge release/v1.1.0
git push origin main
```

---

## For Contributors

### Version Info in Development

When developing, your version will show commit hash:

```bash
$ make build
$ ./bin/gohome --version
gohome abc123f-dirty (commit: abc123f, built: 2026-01-08T...)
```

The `-dirty` suffix indicates **uncommitted changes**. Commit before testing:

```bash
git add .
git commit -m "feat: my awesome feature"
make build
./bin/gohome --version
# gohome abc123f (commit: abc123f, ...) ← no more -dirty
```

### Testing Version Detection

Test different scenarios:

```bash
# 1. Local build with ldflags (via Makefile)
make build
./bin/gohome --version

# 2. Simulate go install (no ldflags)
go build -o /tmp/gohome ./cmd/gohome
/tmp/gohome --version

# 3. Check VCS info available
go build -o /tmp/gohome ./cmd/gohome
go version -m /tmp/gohome | grep vcs
```

---

## Troubleshooting

### Issue: Version shows "dev" always

**Symptoms:**
```bash
$ gohome --version
gohome dev (commit: none, built: unknown)
```

**Causes & Solutions:**

1. **Built without Makefile**
   ```bash
   # ❌ Wrong
   go build ./cmd/gohome
   
   # ✅ Correct
   make build
   ```

2. **Git not installed or not in git repo**
   ```bash
   # Check
   git describe --tags --always
   
   # If error, you're not in a git repository
   git init
   git add .
   git commit -m "Initial commit"
   ```

3. **No git tags exist**
   ```bash
   # Check tags
   git tag -l
   
   # If empty, create one
   git tag v0.1.0
   ```

### Issue: Go install shows pseudo-version

**Symptoms:**
```bash
$ go install github.com/anIcedAntFA/gohome/cmd/gohome@latest
$ gohome --version
gohome v0.0.0-20260108065610-abc123 (commit: none, built: unknown)
```

**This is expected!** Go install uses VCS pseudo-versions when no stable release exists.

**Solution:**
- Create an official release: `git tag v1.0.0 && git push origin v1.0.0`
- Install specific version: `go install ...@v1.0.0`

### Issue: Wrong version in GitHub Release binaries

**Symptoms:**
Binary shows wrong version or "dev"

**Causes:**
- `.goreleaser.yml` has incorrect ldflags path
- Build happened from wrong commit/tag

**Solution:**
1. Check `.goreleaser.yml`:
   ```yaml
   ldflags:
     - -X github.com/anIcedAntFA/gohome/internal/version.Version={{.Version}}
     # NOT: -X main.version={{.Version}}  ❌
   ```

2. Re-tag if needed:
   ```bash
   git tag -d v1.0.0           # Delete local tag
   git push origin :v1.0.0     # Delete remote tag
   git tag v1.0.0              # Re-create tag
   git push origin v1.0.0      # Push again
   ```

### Issue: Version in CI differs from local

**Symptoms:**
- Local: `v1.0.0`
- CI: `abc123-dirty`

**Causes:**
- CI not checking out full git history
- CI not fetching tags

**Solution:**

Check `.github/workflows/*.yml`:
```yaml
- name: Checkout code
  uses: actions/checkout@v4
  with:
    fetch-depth: 0  # ← Important! Fetch all history and tags
```

---

## Additional Resources

- [Semantic Versioning](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Go Modules Version Numbers](https://go.dev/ref/mod#versions)
- [Git Tagging](https://git-scm.com/book/en/v2/Git-Basics-Tagging)

---

**Last Updated:** 2026-01-08
**Version:** 1.0.0
