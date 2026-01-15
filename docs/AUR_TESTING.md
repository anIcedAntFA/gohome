# Testing AUR Configuration Locally

This guide helps you test the AUR configuration before pushing to production.

## Prerequisites

Install GoReleaser:

```bash
# Using Go
go install github.com/goreleaser/goreleaser/v2@latest

# Or using Homebrew (macOS/Linux)
brew install goreleaser/tap/goreleaser

# Or download binary from https://github.com/goreleaser/goreleaser/releases
```

## Test Build (Without Publishing)

### 1. Build Snapshot

Build without creating a release or pushing to AUR:

```bash
cd /home/ngockhoi96/workspace/github.com/anIcedAntFA/gohome

# Build for your current platform only
goreleaser build --snapshot --clean --single-target

# Check the output
ls -la dist/
```

### 2. Generate Full Release (Local Only)

Test the complete release process without publishing:

```bash
# This will generate everything including PKGBUILD
goreleaser release --snapshot --skip=publish --clean

# Check generated AUR files
ls -la dist/
cat dist/aur_sources/gohome/PKGBUILD
cat dist/aur_sources/gohome/.SRCINFO
```

Expected output structure:
```
dist/
├── aur_sources/
│   └── gohome/
│       ├── PKGBUILD
│       └── .SRCINFO
├── checksums.txt
├── gohome_*_linux_x86_64.tar.gz
├── gohome_*_darwin_arm64.tar.gz
└── ...
```

### 3. Validate PKGBUILD Syntax

If you're on Arch Linux or have `namcap` installed:

```bash
cd dist/aur_sources/gohome

# Check PKGBUILD for errors
namcap PKGBUILD

# Build locally (requires Arch Linux)
makepkg -si
```

### 4. Test Installation (Arch Linux Only)

If you have Arch Linux (VM, Docker, or native):

```bash
# Build the package
cd dist/aur_sources/gohome
makepkg -f

# Install
sudo pacman -U gohome-*.pkg.tar.zst

# Test
gohome --version
```

## Using Docker for Testing (Recommended)

If you don't have Arch Linux, use Docker:

```bash
# Pull Arch Linux image
docker pull archlinux:latest

# Run container with dist/ mounted
docker run -it -v $(pwd)/dist/aur_sources/gohome:/pkg archlinux:latest bash

# Inside container:
cd /pkg
pacman -Syu --noconfirm base-devel git
makepkg -si --noconfirm

# Test
gohome --version
```

## Common Issues

### "cannot find package"

Make sure you ran `goreleaser release --snapshot --skip=publish --clean` first.

### "sha256sums mismatch"

This is expected in snapshot mode. Real releases will have correct checksums.

### "unknown variable 'pkgname'"

Check PKGBUILD syntax. Make sure `.goreleaser.yml` has correct template syntax.

### Build fails with "permission denied"

The build instructions in `.goreleaser.yml` might need adjustment. Check:
- File permissions
- Build flags
- Ldflags syntax

## Next Steps

Once local testing passes:

1. Complete [docs/AUR_SETUP.md](AUR_SETUP.md) setup
2. Commit changes and create a new release tag
3. Monitor GitHub Actions for AUR publish job
4. Verify package appears at https://aur.archlinux.org/packages/gohome

## Useful Commands

```bash
# Validate goreleaser config
goreleaser check

# List what would be built
goreleaser build --snapshot --skip=publish --clean --dry-run

# Build for specific OS
GOOS=linux GOARCH=amd64 goreleaser build --snapshot --single-target

# Clean dist/
rm -rf dist/
```
