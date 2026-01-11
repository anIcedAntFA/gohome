# NPM Publishing Guide

This guide explains how gohome is published to npm as a wrapper package for the Go binary.

## Architecture

The npm package (`@ngockhoi96/gohome`) is a **binary wrapper** that:
1. Downloads the pre-compiled binary from GitHub Releases
2. Extracts it to `node_modules/@anicedantfa/gohome/bin/`
3. Exposes it as a global command when installed with `-g`

**Why this approach?**
- GoReleaser's `npms` feature is Pro-only (requires paid license)
- Manual wrapper provides same user experience with open-source tools
- Leverages existing GitHub Releases infrastructure

## Package Structure

```
npm/
├── package.json      # NPM package metadata
├── install.js        # Post-install script (downloads binary)
├── README.md         # NPM-specific documentation
└── .npmignore        # Exclude unnecessary files from package
```

## How It Works

### 1. User Installs Package
```bash
npm install -g @ngockhoi96/gohome
```

### 2. Post-Install Script Runs
The `install.js` script:
- Detects user's platform (darwin/linux/windows) and architecture (amd64/arm64)
- Constructs GitHub Release download URL
- Downloads the appropriate binary archive
- Extracts binary to `bin/` directory
- Makes it executable (on Unix systems)

### 3. Binary Available Globally
The `bin` field in `package.json` creates symlink:
```json
{
  "bin": {
    "gohome": "bin/gohome"
  }
}
```

## Publishing Process

### Automatic (via GitHub Actions)

**Trigger:** Push git tag (e.g., `v1.0.3`)

**Workflow:** `.github/workflows/release.yml`

1. **`release` job:**
   - Builds cross-platform binaries
   - Creates GitHub Release
   - Uploads binaries as release assets

2. **`publish-npm` job:** (runs after `release`)
   - Extracts version from git tag
   - Updates `npm/package.json` version
   - Publishes to npm registry

**Requirements:**
- `NPM_TOKEN` secret must be set in repository settings
- Token must have "Automation" provenance setting

### Manual Publishing

1. **Update version in package.json:**
   ```bash
   cd npm
   npm version 1.0.3 --no-git-tag-version
   ```

2. **Test locally:**
   ```bash
   npm install -g .
   gohome --version
   ```

3. **Publish to npm:**
   ```bash
   npm publish --access public
   ```

## NPM Token Setup

### Creating Token (npmjs.com)

1. Go to https://www.npmjs.com/settings/YOUR_USERNAME/tokens
2. Click "Generate New Token" → "Classic Token"
3. Select "Automation" (for CI/CD)
4. Copy token

### Adding Token to GitHub

1. Go to repository → Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `NPM_TOKEN`
4. Value: (paste token)
5. Click "Add secret"

## Package Scope

**Package name:** `@ngockhoi96/gohome`

The `@ngockhoi96` scope:
- Already exists with your npm account (https://www.npmjs.com/~ngockhoi96)
- User scope (automatic with your username)
- No organization needed
- Just need npm token for publishing

## Supported Platforms

| Platform | Architectures | Archive Format |
|----------|---------------|----------------|
| Linux    | amd64, arm64  | tar.gz         |
| macOS    | amd64, arm64  | tar.gz         |
| Windows  | amd64, arm64  | zip            |

These match GoReleaser's build matrix in `.goreleaser.yml`.

## Testing NPM Package Locally

### Test installation from GitHub Release

1. **Ensure release exists:**
   ```bash
   # Check if release v1.0.2 exists with binaries
   curl -I https://github.com/anIcedAntFA/gohome/releases/download/v1.0.2/gohome_1.0.2_linux_amd64.tar.gz
   ```

2. **Test install script:**
   ```bash
   cd npm
   node install.js
   ./bin/gohome --version
   ```

3. **Test as npm package:**
   ```bash
   npm install -g .
   gohome --version
   ```

4. **Test npx (without installation):**
   ```bash
   npx /path/to/gohome/npm --help
   ```

### Test with specific version

```bash
# Test downloading specific release version
cd npm
npm version 1.0.2 --no-git-tag-version
node install.js
```

## Troubleshooting

### Install Script Fails

**Error:** `Download failed with status 404`
- **Cause:** GitHub Release doesn't exist for specified version
- **Fix:** Ensure `package.json` version matches existing GitHub Release tag

**Error:** `Unsupported platform`
- **Cause:** Platform/architecture not supported
- **Fix:** Check `PLATFORM_MAPPING` and `ARCH_MAPPING` in `install.js`

**Error:** `Extraction failed`
- **Cause:** Missing tar/unzip utilities
- **Fix:** Install required tools:
  - Linux: `sudo apt-get install tar`
  - Windows: Built-in `Expand-Archive`

### Publish Fails

**Error:** `403 Forbidden`
- **Cause:** Insufficient npm token permissions
- **Fix:** Regenerate token with "Automation" scope

**Error:** `404 Not Found - PUT https://registry.npmjs.org/@ngockhoi96%2fgohome`
- **Cause:** Scope `@ngockhoi96` doesn't exist (unlikely if using your account)
- **Fix:** Verify you're logged in: `npm whoami`

**Error:** `Package name too similar to existing packages`
- **Cause:** npm name squatting protection
- **Fix:** Choose different package name

### Version Mismatch

**Error:** Binary version doesn't match package version
- **Cause:** `npm/package.json` not updated during release
- **Fix:** GitHub Actions workflow extracts version from git tag automatically

## Maintenance

### Update Dependencies

The package has minimal dependencies (Node.js built-ins only):
- `fs` - File system operations
- `path` - Path manipulation
- `https` - Binary download
- `child_process` - Archive extraction

**No npm dependencies to update!**

### Version Synchronization

**Critical:** Keep versions synchronized:
- Git tag: `v1.0.3`
- `npm/package.json`: `"version": "1.0.3"`
- GitHub Release: `v1.0.3`

**Automation:** GitHub Actions workflow automatically:
1. Extracts version from git tag (`v1.0.3` → `1.0.3`)
2. Updates `package.json` with `npm version`
3. Publishes to npm

## Migration Notes

### From GoReleaser Pro `npms` Section

If upgrading to GoReleaser Pro in the future:

1. **Remove manual npm structure:**
   ```bash
   rm -rf npm/
   ```

2. **Add `npms` to `.goreleaser.yml`:**
   ```yaml
   npms:
     - name: "@anicedantfa/gohome"
       access: public
       format: tgz
       # ... (see GoReleaser Pro docs)
   ```

3. **Simplify GitHub Actions:**
   - Remove separate `publish-npm` job
   - GoReleaser handles publishing automatically

**Trade-offs:**
- Pro: Integrated with GoReleaser, less maintenance
- Con: Requires paid license ($99/year as of 2025)

## References

- [npm documentation](https://docs.npmjs.com/)
- [npm binary wrapper pattern](https://docs.npmjs.com/cli/v8/configuring-npm/package-json#bin)
- [GoReleaser npms (Pro)](https://goreleaser.com/customization/npm/)
- [GitHub Actions Node.js setup](https://github.com/actions/setup-node)

## Support

Questions or issues with npm package?
- Open issue: https://github.com/anIcedAntFA/gohome/issues
- Check existing: https://github.com/anIcedAntFA/gohome/issues?q=is%3Aissue+label%3Anpm
