# 🚀 Quick Start Checklist

Fast track checklist to get your first release out.

## Phase 1: Minimal Setup (30 minutes)

### ✅ Step 1: Add Version Support (5 min)

- [ ] Update `cmd/gohome/main.go`:
  ```go
  var (
      version = "dev"
      commit  = "none"
      date    = "unknown"
  )
  
  func main() {
      showVersion := flag.Bool("version", false, "Show version")
      flag.BoolVar(showVersion, "v", false, "Show version (shorthand)")
      // ... existing flags
      flag.Parse()
      
      if *showVersion {
          fmt.Printf("gohome %s (commit: %s, built: %s)\n", version, commit, date)
          os.Exit(0)
      }
      // ... rest of code
  }
  ```

### ✅ Step 2: Create GoReleaser Config (10 min)

- [ ] Install GoReleaser:
  ```bash
  brew install goreleaser/tap/goreleaser
  # or
  go install github.com/goreleaser/goreleaser@latest
  ```

- [ ] Create `.goreleaser.yaml` (copy from RELEASE_GUIDE.md section "GoReleaser Setup")
- [ ] Update your info:
  - `project_name`
  - `maintainer`
  - `homepage`
  - GitHub username/repo

- [ ] Test config:
  ```bash
  goreleaser check
  goreleaser build --snapshot --clean
  ```

### ✅ Step 3: Create Installation Script (5 min)

- [ ] Copy `install.sh` from RELEASE_GUIDE.md
- [ ] Update `REPO="anIcedAntFA/gohome"` with your username
- [ ] Make executable:
  ```bash
  chmod +x install.sh
  ```

### ✅ Step 4: Setup GitHub Actions (10 min)

- [ ] Create `.github/workflows/release.yml` (copy from guide)
- [ ] Verify GITHUB_TOKEN permissions:
  - Go to: Settings → Actions → General → Workflow permissions
  - Enable: "Read and write permissions"

- [ ] Commit and push:
  ```bash
  git add .goreleaser.yaml install.sh .github/workflows/release.yml
  git commit -m "chore: setup release automation"
  git push
  ```

## Phase 2: First Release (10 minutes)

### ✅ Step 5: Create Release

- [ ] Update CHANGELOG.md:
  ```markdown
  ## v1.0.0 - 2026-01-07
  
  ### Features
  - Initial release
  - Git commit aggregation
  - Custom tasks support
  - Multiple output formats
  ```

- [ ] Commit:
  ```bash
  git add CHANGELOG.md
  git commit -m "chore: prepare v1.0.0 release"
  git push
  ```

- [ ] Create tag:
  ```bash
  git tag -a v1.0.0 -m "Release v1.0.0"
  git push origin v1.0.0
  ```

- [ ] Watch GitHub Actions:
  - Go to: https://github.com/YOUR_USERNAME/gohome/actions
  - Wait for "Release" workflow to complete (~5 min)

- [ ] Verify release:
  - Go to: https://github.com/YOUR_USERNAME/gohome/releases
  - Check assets are uploaded

### ✅ Step 6: Test Installation

- [ ] Test `go install`:
  ```bash
  go install github.com/YOUR_USERNAME/gohome/cmd/gohome@v1.0.0
  gohome --version
  ```

- [ ] Test `install.sh`:
  ```bash
  curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/gohome/main/install.sh | bash
  gohome --version
  ```

## Phase 3: Enhanced Distribution (Optional, 1-2 hours)

### ✅ Step 7: Homebrew Tap

- [ ] Create new repo: `homebrew-tap`
  ```bash
  gh repo create homebrew-tap --public
  ```

- [ ] Update `.goreleaser.yaml` brews section with your username
- [ ] Next release will auto-push formula

- [ ] Test:
  ```bash
  brew tap YOUR_USERNAME/tap
  brew install gohome
  ```

### ✅ Step 8: Update README

- [ ] Add installation section:
  ```markdown
  ## Installation
  
  ### macOS
  ```bash
  brew install YOUR_USERNAME/tap/gohome
  ```
  
  ### Linux/macOS
  ```bash
  curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/gohome/main/install.sh | bash
  ```
  
  ### Go Install
  ```bash
  go install github.com/YOUR_USERNAME/gohome/cmd/gohome@latest
  ```
  
  ### Download Binary
  Download from [releases page](https://github.com/YOUR_USERNAME/gohome/releases/latest)
  ```

- [ ] Add badge:
  ```markdown
  ![Release](https://img.shields.io/github/v/release/YOUR_USERNAME/gohome)
  ```

### ✅ Step 9: npm Package (Optional)

- [ ] Create `package.json` (see RELEASE_GUIDE.md)
- [ ] Create download script
- [ ] Publish to npm:
  ```bash
  npm login
  npm publish --access public
  ```

## Quick Commands Reference

```bash
# Check GoReleaser config
goreleaser check

# Build locally without releasing
goreleaser build --snapshot --clean

# Test with specific version
goreleaser build --snapshot --clean --single-target

# Create release tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Delete tag if needed
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# View GitHub release
gh release view v1.0.0

# Download release asset
gh release download v1.0.0
```

## Troubleshooting

### GoReleaser fails with "invalid config"
```bash
# Check YAML syntax
goreleaser check

# Validate structure
cat .goreleaser.yaml | yq
```

### GitHub Actions fails
```bash
# Check logs
gh run list --workflow=release.yml
gh run view [RUN_ID] --log

# Re-run failed jobs
gh run rerun [RUN_ID]
```

### Permission denied on install.sh
```bash
# Make executable
chmod +x install.sh

# Or run with bash explicitly
bash install.sh
```

## Success Criteria

After completing these steps, you should have:

- ✅ `gohome --version` shows correct version
- ✅ GitHub Releases page has binaries for all platforms
- ✅ Installation script works on Linux/macOS
- ✅ `go install` works
- ✅ Automated releases via git tags

## Next Release

For subsequent releases:

```bash
# 1. Make changes
# 2. Update CHANGELOG.md
# 3. Commit
git add .
git commit -m "feat: new feature"
git push

# 4. Create tag
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# Done! GitHub Actions handles the rest
```

---

**Need help?** Check [RELEASE_GUIDE.md](RELEASE_GUIDE.md) for detailed explanations.
