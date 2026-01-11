# NPM Support Implementation Summary

## Overview

Successfully implemented NPM package distribution for gohome using a **binary wrapper pattern**. The package `@anicedantfa/gohome` downloads pre-built binaries from GitHub Releases during installation.

## Why Manual Wrapper Instead of GoReleaser Pro?

GoReleaser's `npms` feature is **Pro-only** (requires $99/year license as of 2025). The manual wrapper approach:
- ✅ Achieves identical user experience
- ✅ Uses only open-source tools
- ✅ Leverages existing GitHub Releases infrastructure
- ✅ Zero additional cost

## Implementation Details

### Package Structure

```
npm/
├── package.json      # Package metadata, version 1.0.2
├── install.js        # Post-install script (downloads binary)
├── README.md         # NPM-specific documentation
└── .npmignore        # Excludes source code from package
```

### How It Works

1. **User installs:** `npm install -g @anicedantfa/gohome`
2. **Post-install runs:** `install.js` detects platform/arch, downloads binary from GitHub Release
3. **Binary extracted:** To `node_modules/@anicedantfa/gohome/bin/gohome`
4. **Symlink created:** npm links `bin/gohome` to global PATH

### Supported Platforms

| Platform | Architectures | Archive Format |
|----------|---------------|----------------|
| Linux    | amd64, arm64  | tar.gz         |
| macOS    | amd64, arm64  | tar.gz         |
| Windows  | amd64, arm64  | zip            |

## CI/CD Integration

### GitHub Actions Workflow

`.github/workflows/release.yml` now has two jobs:

1. **`release` job:**
   - Builds cross-platform binaries (GoReleaser)
   - Creates GitHub Release with binaries

2. **`publish-npm` job:** (new)
   - Runs after `release` job completes
   - Extracts version from git tag (`v1.0.3` → `1.0.3`)
   - Updates `npm/package.json` version
   - Publishes to npmjs.com registry

### Publishing Workflow

```
Push tag (v1.0.3) → GoReleaser builds → GitHub Release → NPM publish
```

**Automatic:** Triggered by git tag push
**Manual:** `cd npm && npm publish --access public`

## Documentation

### User-Facing

- **README.md:** Added NPM installation section
  ```bash
  # Install globally
  npm install -g @anicedantfa/gohome
  
  # Or use without installation
  npx @anicedantfa/gohome --help
  ```

- **CHANGELOG.md:** Documented in [Unreleased] section
- **ROADMAP.md:** Marked npm as completed ✅
- **.goreleaser.yml:** Added npm to release notes

### Developer Documentation

- **docs/NPM_PUBLISHING.md:** Comprehensive 271-line guide
  - Architecture explanation
  - Manual/automatic publishing
  - Token setup instructions
  - Platform support details
  - Troubleshooting guide

## Testing Results

✅ **goreleaser check:** Configuration valid  
✅ **npm pack --dry-run:** Package structure verified (3 files, 5.9 kB unpacked)  
✅ **node -c install.js:** Syntax validated

## Setup Requirements (Before First Publish)

### 1. Create npm Account & Scope

```bash
# Create account at https://www.npmjs.com/signup
# Create scope @anicedantfa (organization or user scope)
```

### 2. Generate NPM Token

1. Go to https://www.npmjs.com/settings/YOUR_USERNAME/tokens
2. Click "Generate New Token" → "Classic Token"
3. Select **"Automation"** (for CI/CD)
4. Copy token

### 3. Add Token to GitHub

1. Go to repository → Settings → Secrets and variables → Actions
2. New repository secret
3. Name: `NPM_TOKEN`
4. Value: (paste token)

### 4. Test Publish

With next release (v1.0.3+):
- Push tag: `git tag v1.0.3 && git push origin v1.0.3`
- GitHub Actions automatically publishes
- Verify: `npm view @anicedantfa/gohome`

## Files Changed

### New Files (7)
```
npm/package.json         - NPM package metadata
npm/install.js           - Binary download/extraction script
npm/README.md            - NPM-specific docs
npm/.npmignore           - Package exclusions
docs/NPM_PUBLISHING.md   - Publishing guide
```

### Modified Files (6)
```
.github/workflows/release.yml  - Added publish-npm job
.gitignore                     - Exclude npm/bin/ and *.tgz
.goreleaser.yml                - Added npm to release notes
CHANGELOG.md                   - Document NPM support
README.md                      - Installation section
ROADMAP.md                     - Mark npm completed
```

### Deleted Files (2)
```
docs/demos/demo-quickstart.sh  - Old bash demo script
docs/demos/demo-config.sh      - Old bash demo script
```

## Commit Details

**Commit:** `5061979`  
**Branch:** `feat/support-npm`  
**Changes:** +600 insertions, -112 deletions  
**Message:** ✨ feat(distribution): add NPM package support

## Usage Examples

### Install Globally
```bash
npm install -g @anicedantfa/gohome
gohome --help
```

### Use with npx (No Installation)
```bash
npx @anicedantfa/gohome --days 3
npx @anicedantfa/gohome --copy
```

### Uninstall
```bash
npm uninstall -g @anicedantfa/gohome
```

## Advantages of This Approach

1. **Zero Dependencies:** Only Node.js built-ins (fs, path, https, child_process)
2. **Small Package:** 5.9 kB unpacked (binaries downloaded separately)
3. **Fast Downloads:** Downloads only platform-specific binary (not all platforms)
4. **Offline Friendly:** Works with npm cache, binaries cached by GitHub CDN
5. **Version Sync:** Automated via GitHub Actions (tag → package.json)

## Limitations

1. **Requires npm ≥14:** Uses modern Node.js APIs
2. **Network Required:** Post-install downloads from GitHub
3. **GitHub Dependency:** Relies on GitHub Releases availability
4. **No Scoped Permissions:** Token has full account access (use least-privilege)

## Future Improvements

- [ ] Add checksums verification in install.js (SHA256)
- [ ] Support proxies (HTTP_PROXY environment variable)
- [ ] Add retry logic for download failures
- [ ] Cache downloaded binaries (avoid re-download on reinstall)
- [ ] Migrate to GoReleaser Pro `npms` if/when license acquired

## Related Resources

- [npm binary wrapper pattern](https://docs.npmjs.com/cli/v8/configuring-npm/package-json#bin)
- [GoReleaser npms (Pro)](https://goreleaser.com/customization/npm/)
- [GitHub Actions Node.js setup](https://github.com/actions/setup-node)
- [Publishing guide](docs/NPM_PUBLISHING.md)

## Support

Questions or issues:
- Open issue: https://github.com/anIcedAntFA/gohome/issues
- Label: `npm` for npm-specific issues

---

**Status:** ✅ Implementation complete, ready for first publish with v1.0.3+
