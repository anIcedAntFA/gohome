# 🚀 Release & Distribution Guide

Complete guide for building, releasing, and distributing **gohome** across multiple platforms and package managers.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Versioning Strategy](#versioning-strategy)
4. [GoReleaser Setup](#goreleaser-setup)
5. [GitHub Actions CI/CD](#github-actions-cicd)
6. [Distribution Channels](#distribution-channels)
7. [Installation Methods](#installation-methods)
8. [Testing Releases](#testing-releases)
9. [Troubleshooting](#troubleshooting)

---

## Overview

### Distribution Matrix

| Platform | Architecture | Package Format | Installation Method |
|----------|-------------|----------------|-------------------|
| **Linux** | amd64, arm64 | `.tar.gz`, `.deb`, `.rpm` | `go install`, `curl\|sh`, apt, yum, snap |
| **macOS** | amd64 (Intel), arm64 (M1/M2) | `.tar.gz`, `.dmg` | `go install`, `curl\|sh`, Homebrew |
| **Windows** | amd64, arm64 | `.zip`, `.exe` | `go install`, Scoop, Chocolatey |

### Release Automation Flow

```
┌─────────────┐
│ Git Tag Push│
│   v1.0.0    │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│ GitHub Actions  │
│   Triggered     │
└──────┬──────────┘
       │
       ▼
┌─────────────────┐
│  GoReleaser     │
│  - Build        │
│  - Archive      │
│  - Checksum     │
└──────┬──────────┘
       │
       ├─────────────────────┬─────────────────┬──────────────┐
       ▼                     ▼                 ▼              ▼
┌─────────────┐    ┌──────────────┐   ┌────────────┐  ┌──────────┐
│GitHub Release│    │ Homebrew Tap │   │ Snapcraft  │  │  npm     │
│  + Assets   │    │   Formula    │   │   Store    │  │ Package  │
└─────────────┘    └──────────────┘   └────────────┘  └──────────┘
```

---

## Prerequisites

### 1. Install Required Tools

```bash
# GoReleaser
brew install goreleaser/tap/goreleaser
# or
go install github.com/goreleaser/goreleaser@latest

# GitHub CLI (for authentication)
brew install gh
# or download from https://cli.github.com
```

### 2. GitHub Token Setup

```bash
# Create a Personal Access Token (Classic)
# Go to: https://github.com/settings/tokens/new
# Scopes needed: repo, write:packages

# Set in your environment
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"

# Or authenticate with gh CLI
gh auth login
```

### 3. Project Structure

Ensure your project has this structure:

```
gohome/
├── cmd/
│   └── gohome/
│       └── main.go
├── internal/
├── .goreleaser.yaml       # ← We'll create this
├── .github/
│   └── workflows/
│       └── release.yml    # ← We'll create this
├── install.sh             # ← Installation script
├── go.mod
└── README.md
```

---

## Versioning Strategy

### Semantic Versioning (SemVer)

Follow `MAJOR.MINOR.PATCH` format:

- **MAJOR:** Breaking changes (e.g., `v1.0.0` → `v2.0.0`)
- **MINOR:** New features, backward compatible (e.g., `v1.0.0` → `v1.1.0`)
- **PATCH:** Bug fixes (e.g., `v1.0.0` → `v1.0.1`)

### Version Injection

**Update `main.go` to support `--version` flag:**

```go
// cmd/gohome/main.go
package main

import (
    "flag"
    "fmt"
    "os"
    // ... other imports
)

// These will be set by GoReleaser during build
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // Add version flag
    showVersion := flag.Bool("version", false, "Show version information")
    flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")
    
    // ... other flags
    
    flag.Parse()
    
    if *showVersion {
        fmt.Printf("gohome %s (commit: %s, built: %s)\n", version, commit, date)
        os.Exit(0)
    }
    
    // ... rest of main logic
}
```

### Creating Releases

```bash
# 1. Update CHANGELOG.md with new features/fixes

# 2. Commit changes
git add .
git commit -m "chore: prepare release v1.0.0"

# 3. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. GitHub Actions will automatically build and release
```

---

## GoReleaser Setup

### 1. Initialize GoReleaser Config

```bash
cd /path/to/gohome
goreleaser init
```

### 2. Create `.goreleaser.yaml`

Create a comprehensive config file:

```yaml
# .goreleaser.yaml
version: 2

# Project metadata
project_name: gohome

before:
  hooks:
    # Ensure dependencies are up to date
    - go mod tidy
    - go mod download

# Build configuration
builds:
  - id: gohome
    main: ./cmd/gohome
    binary: gohome
    
    # Version injection
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
    
    # Target platforms
    goos:
      - linux
      - darwin  # macOS
      - windows
    
    goarch:
      - amd64
      - arm64
      - arm     # For Raspberry Pi
    
    # Ignore unsupported combinations
    ignore:
      - goos: darwin
        goarch: arm  # macOS doesn't run on 32-bit ARM
      - goos: windows
        goarch: arm  # Windows ARM is arm64 only

# Archives for different platforms
archives:
  - id: default
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    
    format_overrides:
      - goos: windows
        format: zip
    
    files:
      - README.md
      - LICENSE
      - CHANGELOG.md
      - install.sh  # Include installation script

# Linux packages (deb, rpm)
nfpms:
  - id: gohome
    package_name: gohome
    
    vendor: anIcedAntFA
    homepage: https://github.com/anIcedAntFA/gohome
    maintainer: Your Name <your.email@example.com>
    
    description: |
      A CLI tool for generating daily standup reports from git commits.
      Supports multiple repositories, custom tasks, and various output formats.
    
    license: MIT
    
    formats:
      - deb   # Debian/Ubuntu
      - rpm   # RedHat/Fedora/CentOS
      - apk   # Alpine Linux
    
    # Package dependencies
    dependencies:
      - git
    
    # Installation scripts
    scripts:
      postinstall: "scripts/postinstall.sh"
    
    # File mappings
    contents:
      - src: ./gohome
        dst: /usr/local/bin/gohome
      - src: ./README.md
        dst: /usr/share/doc/gohome/README.md
      - src: ./LICENSE
        dst: /usr/share/doc/gohome/LICENSE

# Homebrew tap (for macOS)
brews:
  - name: gohome
    
    # Repository to push formula
    repository:
      owner: anIcedAntFA
      name: homebrew-tap
      branch: main
      token: "{{ .Env.GITHUB_TOKEN }}"
    
    # Commit message
    commit_author:
      name: goreleaserbot
      email: bot@goreleaser.com
    
    # Formula details
    homepage: https://github.com/anIcedAntFA/gohome
    description: "Daily standup report generator from git commits"
    license: MIT
    
    # Installation
    install: |
      bin.install "gohome"
    
    # Test
    test: |
      system "#{bin}/gohome", "--version"

# Snapcraft (Linux universal package)
snapcrafts:
  - name: gohome
    
    summary: Daily standup report generator
    description: |
      Generate daily standup reports from git commits.
      Supports multiple repositories and custom tasks.
    
    grade: stable
    confinement: strict
    
    publish: true
    
    apps:
      gohome:
        command: gohome
        plugs: ["home", "network"]

# Checksum of all artifacts
checksum:
  name_template: 'checksums.txt'
  algorithm: sha256

# Source code archive
source:
  enabled: true

# Changelog generation
changelog:
  use: github
  sort: asc
  
  groups:
    - title: Features
      regexp: '^.*?feat(\([[:word:]]+\))??!?:.+$'
      order: 0
    - title: Bug Fixes
      regexp: '^.*?fix(\([[:word:]]+\))??!?:.+$'
      order: 1
    - title: Documentation
      regexp: '^.*?docs(\([[:word:]]+\))??!?:.+$'
      order: 2
    - title: Others
      order: 999
  
  filters:
    exclude:
      - '^chore:'
      - '^test:'
      - '^ci:'

# GitHub Release
release:
  github:
    owner: anIcedAntFA
    name: gohome
  
  draft: false
  prerelease: auto
  mode: replace
  
  # Release notes template
  header: |
    ## 🎉 gohome {{ .Tag }}
    
    **Release Date:** {{ .Date }}
    
    ### 📦 Installation
    
    #### macOS
    ```bash
    brew install anIcedAntFA/tap/gohome
    ```
    
    #### Linux
    ```bash
    curl -sSL https://raw.githubusercontent.com/anIcedAntFA/gohome/main/install.sh | bash
    ```
    
    #### Go Install
    ```bash
    go install github.com/anIcedAntFA/gohome/cmd/gohome@{{ .Tag }}
    ```
    
    #### Direct Download
    Download binaries from the assets below.
  
  footer: |
    ---
    **Full Changelog**: https://github.com/anIcedAntFA/gohome/compare/{{ .PreviousTag }}...{{ .Tag }}

# Announce release (optional)
announce:
  skip: "{{gt .Patch 0}}"  # Skip patch releases
  
  # Twitter
  twitter:
    enabled: false
    message_template: 'gohome {{ .Tag }} is out! Check it out at {{ .ReleaseURL }}'
```

---

## GitHub Actions CI/CD

### 1. Create Release Workflow

Create `.github/workflows/release.yml`:

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write
  packages: write

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for changelog
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true
      
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.txt ./...
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          # Add other secrets if needed (e.g., for Homebrew tap)
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.txt
```

### 2. Create Test & Lint Workflow

Update or create `.github/workflows/ci.yml`:

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: v2.7.2
  
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.txt
  
  build:
    name: Build
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.23']
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      
      - name: Build
        run: go build -v ./cmd/gohome
```

---

## Distribution Channels

### 1. Direct Binary Downloads (GitHub Releases)

**Automatic via GoReleaser** ✅

Users can download from: `https://github.com/anIcedAntFA/gohome/releases/latest`

### 2. Installation Script (`install.sh`)

Create a universal installation script:

```bash
#!/bin/bash
# install.sh - Universal installer for gohome

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Constants
REPO="anIcedAntFA/gohome"
BINARY_NAME="gohome"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="x86_64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="armv7"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${BLUE}Detected platform: $OS/$ARCH${NC}"
}

# Get latest version
get_latest_version() {
    echo -e "${BLUE}Fetching latest version...${NC}"
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Failed to fetch latest version${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Latest version: $VERSION${NC}"
}

# Download and install
install_binary() {
    local filename="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
    local download_url="https://github.com/$REPO/releases/download/$VERSION/$filename"
    local tmp_dir=$(mktemp -d)
    
    echo -e "${BLUE}Downloading $filename...${NC}"
    
    if ! curl -sL "$download_url" -o "$tmp_dir/$filename"; then
        echo -e "${RED}Download failed${NC}"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    echo -e "${BLUE}Extracting...${NC}"
    tar -xzf "$tmp_dir/$filename" -C "$tmp_dir"
    
    echo -e "${BLUE}Installing to $INSTALL_DIR...${NC}"
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        echo -e "${YELLOW}Requesting sudo access to install to $INSTALL_DIR${NC}"
        sudo mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    rm -rf "$tmp_dir"
    
    echo -e "${GREEN}✓ Successfully installed $BINARY_NAME $VERSION${NC}"
    echo -e "${GREEN}✓ Run '$BINARY_NAME --version' to verify${NC}"
}

# Main
main() {
    echo -e "${BLUE}═══════════════════════════════════════${NC}"
    echo -e "${BLUE}   gohome Installation Script${NC}"
    echo -e "${BLUE}═══════════════════════════════════════${NC}\n"
    
    detect_platform
    get_latest_version
    install_binary
    
    echo -e "\n${BLUE}═══════════════════════════════════════${NC}"
    echo -e "${GREEN}Installation complete!${NC}"
    echo -e "${BLUE}═══════════════════════════════════════${NC}\n"
    
    # Show usage hint
    echo -e "${YELLOW}Quick Start:${NC}"
    echo -e "  $BINARY_NAME -d 7          # Show commits from last 7 days"
    echo -e "  $BINARY_NAME --help        # Show all options"
    echo -e "  $BINARY_NAME --version     # Show version\n"
}

main "$@"
```

**Usage:**
```bash
curl -sSL https://raw.githubusercontent.com/anIcedAntFA/gohome/main/install.sh | bash
```

### 3. Homebrew (macOS/Linux)

**Setup Homebrew Tap:**

1. Create a new repository: `homebrew-tap`
2. GoReleaser will automatically push formula updates

**Users install via:**
```bash
brew tap anIcedAntFA/tap
brew install gohome

# Or one-liner
brew install anIcedAntFA/tap/gohome
```

### 4. Snapcraft (Linux Universal)

**Automatic via GoReleaser** ✅

Users install via:
```bash
sudo snap install gohome
```

### 5. Package Managers

#### Debian/Ubuntu (APT)

```bash
# Download .deb from releases
wget https://github.com/anIcedAntFA/gohome/releases/download/v1.0.0/gohome_1.0.0_linux_amd64.deb

# Install
sudo dpkg -i gohome_1.0.0_linux_amd64.deb
```

#### RedHat/Fedora (YUM/DNF)

```bash
# Download .rpm from releases
wget https://github.com/anIcedAntFA/gohome/releases/download/v1.0.0/gohome_1.0.0_linux_amd64.rpm

# Install
sudo rpm -i gohome_1.0.0_linux_amd64.rpm
# or
sudo dnf install gohome_1.0.0_linux_amd64.rpm
```

#### Arch Linux (AUR)

Create `PKGBUILD` for AUR (manual):

```bash
# PKGBUILD
pkgname=gohome
pkgver=1.0.0
pkgrel=1
pkgdesc="Daily standup report generator from git commits"
arch=('x86_64' 'aarch64')
url="https://github.com/anIcedAntFA/gohome"
license=('MIT')
depends=('git')
source=("$pkgname-$pkgver.tar.gz::https://github.com/anIcedAntFA/gohome/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$srcdir/$pkgname-$pkgver"
    go build -o $pkgname cmd/gohome/main.go
}

package() {
    cd "$srcdir/$pkgname-$pkgver"
    install -Dm755 $pkgname "$pkgdir/usr/bin/$pkgname"
}
```

Users install via:
```bash
yay -S gohome
# or
paru -S gohome
```

### 6. npm (JavaScript Wrapper)

For users who prefer npm:

**Create `package.json` wrapper:**

```json
{
  "name": "@anicedantfa/gohome",
  "version": "1.0.0",
  "description": "Daily standup report generator from git commits",
  "bin": {
    "gohome": "./bin/gohome"
  },
  "scripts": {
    "postinstall": "node scripts/download-binary.js"
  },
  "keywords": ["git", "standup", "cli", "productivity"],
  "author": "Your Name",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/anIcedAntFA/gohome"
  }
}
```

**Create download script `scripts/download-binary.js`:**

```javascript
const { execSync } = require('child_process');
const https = require('https');
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');
const tar = require('tar');

const REPO = 'anIcedAntFA/gohome';
const VERSION = require('../package.json').version;

const platform = process.platform === 'win32' ? 'windows' : process.platform;
const arch = process.arch === 'x64' ? 'x86_64' : process.arch;

const filename = `gohome_${VERSION}_${platform}_${arch}.tar.gz`;
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;

console.log(`Downloading gohome v${VERSION} for ${platform}/${arch}...`);

// Download and extract
// ... (implement download logic)
```

Users install via:
```bash
npm install -g @anicedantfa/gohome
```

### 7. Windows Package Managers

#### Scoop

Create manifest in `scoop-bucket`:

```json
{
    "version": "1.0.0",
    "description": "Daily standup report generator",
    "homepage": "https://github.com/anIcedAntFA/gohome",
    "license": "MIT",
    "architecture": {
        "64bit": {
            "url": "https://github.com/anIcedAntFA/gohome/releases/download/v1.0.0/gohome_1.0.0_windows_x86_64.zip",
            "hash": "..."
        }
    },
    "bin": "gohome.exe",
    "checkver": "github",
    "autoupdate": {
        "architecture": {
            "64bit": {
                "url": "https://github.com/anIcedAntFA/gohome/releases/download/v$version/gohome_$version_windows_x86_64.zip"
            }
        }
    }
}
```

Users install via:
```powershell
scoop bucket add anicedantfa https://github.com/anIcedAntFA/scoop-bucket
scoop install gohome
```

#### Chocolatey

Create `.nuspec` file (advanced, manual setup required)

---

## Installation Methods

### Summary Table

| Method | Platform | Command | Auto-Update |
|--------|----------|---------|-------------|
| **go install** | All | `go install github.com/anIcedAntFA/gohome/cmd/gohome@latest` | ❌ Manual |
| **curl\|sh** | Unix | `curl -sSL https://raw.githubusercontent.com/.../install.sh \| bash` | ❌ Re-run script |
| **Homebrew** | macOS/Linux | `brew install anIcedAntFA/tap/gohome` | ✅ `brew upgrade` |
| **Snap** | Linux | `sudo snap install gohome` | ✅ Automatic |
| **APT** | Debian/Ubuntu | `sudo dpkg -i gohome.deb` | ❌ Manual download |
| **YUM/DNF** | RedHat/Fedora | `sudo rpm -i gohome.rpm` | ❌ Manual download |
| **AUR** | Arch Linux | `yay -S gohome` | ✅ `yay -Syu` |
| **npm** | All (via Node) | `npm install -g @anicedantfa/gohome` | ✅ `npm update -g` |
| **Scoop** | Windows | `scoop install gohome` | ✅ `scoop update` |
| **Binary** | All | Download from releases | ❌ Manual |

---

## Testing Releases

### 1. Local Test Build

```bash
# Test GoReleaser config
goreleaser check

# Build snapshot (without releasing)
goreleaser build --snapshot --clean

# Check artifacts
ls -la dist/
```

### 2. Test Installation Script

```bash
# Test install.sh locally
bash install.sh

# Verify installation
which gohome
gohome --version
```

### 3. Test GitHub Actions Locally

Use [act](https://github.com/nektos/act):

```bash
# Install act
brew install act

# Test release workflow
act push -j release --secret GITHUB_TOKEN=$GITHUB_TOKEN
```

### 4. Pre-release Testing

Create a pre-release tag:

```bash
git tag -a v1.0.0-rc.1 -m "Release candidate 1"
git push origin v1.0.0-rc.1
```

---

## Step-by-Step Release Process

### First Time Setup

1. **Create necessary repositories:**
   ```bash
   # Create Homebrew tap repo
   gh repo create homebrew-tap --public
   
   # Optional: Create scoop bucket
   gh repo create scoop-bucket --public
   ```

2. **Add GoReleaser config:**
   ```bash
   # Copy .goreleaser.yaml from this guide
   # Customize with your details
   ```

3. **Add GitHub Actions:**
   ```bash
   mkdir -p .github/workflows
   # Copy release.yml and ci.yml
   ```

4. **Add installation script:**
   ```bash
   # Copy install.sh
   chmod +x install.sh
   ```

5. **Update README with installation instructions**

### Regular Release Process

```bash
# 1. Update version in code/docs
# 2. Update CHANGELOG.md
# 3. Commit changes
git add .
git commit -m "chore: bump version to v1.0.0"
git push

# 4. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 5. Wait for GitHub Actions
# Check: https://github.com/anIcedAntFA/gohome/actions

# 6. Verify release
# Check: https://github.com/anIcedAntFA/gohome/releases

# 7. Test installation
brew install anIcedAntFA/tap/gohome
# or
curl -sSL https://raw.githubusercontent.com/.../install.sh | bash

# 8. Announce release (optional)
# - Update README
# - Post on social media
# - Update documentation
```

---

## Troubleshooting

### GoReleaser Fails

```bash
# Check config syntax
goreleaser check

# Build locally to test
goreleaser build --snapshot --clean

# Check logs in GitHub Actions
gh run list --workflow=release.yml
gh run view [RUN_ID] --log
```

### Homebrew Formula Issues

```bash
# Test formula locally
brew install --build-from-source anIcedAntFA/tap/gohome

# Check tap repo
cd $(brew --repository)/Library/Taps/anicedantfa/homebrew-tap
git log
```

### Binary Not Found After Install

```bash
# Check PATH
echo $PATH

# macOS: Add to PATH
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Linux: Usually already in PATH
# Windows: Add to System Environment Variables
```

---

## Best Practices

### 1. Versioning
- Always use semantic versioning
- Tag releases consistently (`v1.0.0`, not `1.0.0`)
- Write meaningful release notes

### 2. Testing
- Test binaries on all platforms before tagging
- Use release candidates for major versions
- Run automated tests in CI before release

### 3. Documentation
- Keep README installation section updated
- Maintain CHANGELOG.md
- Document breaking changes clearly

### 4. Security
- Sign releases (optional but recommended)
- Verify checksums
- Use GitHub's release verification

### 5. User Experience
- Provide multiple installation methods
- Clear error messages in install scripts
- Version compatibility notes

---

## Additional Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Documentation](https://docs.brew.sh/)
- [Snapcraft Documentation](https://snapcraft.io/docs)
- [Semantic Versioning](https://semver.org/)

---

## Next Steps

After completing this guide:

1. ✅ Setup GoReleaser config
2. ✅ Create GitHub Actions workflows
3. ✅ Create installation script
4. ✅ Create first release tag
5. ✅ Test installation methods
6. ✅ Update README with installation instructions
7. ✅ Create Homebrew tap (optional)
8. ✅ Submit to Snapcraft (optional)
9. ✅ Create npm wrapper (optional)

---

**Questions or issues?** Open an issue at https://github.com/anIcedAntFA/gohome/issues
