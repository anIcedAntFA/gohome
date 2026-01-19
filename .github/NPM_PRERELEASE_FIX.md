# NPM Prerelease Publishing Fix

**Date:** 2026-01-19  
**Issue:** NPM publish failed for v1.3.0-beta.2  
**Root Cause:** Missing `--tag` flag for prerelease versions

---

## 🐛 The Problem

### Error Log

```
npm error You must specify a tag using --tag when publishing a prerelease version.
npm error A complete log of this run can be found in: /home/runner/.npm/_logs/2026-01-19T06_34_00_980Z-debug-0.log
Error: Process completed with exit code 1.
```

### Root Cause Analysis

**NPM Behavior:**
- When publishing a version with prerelease suffix (e.g., `-beta.2`, `-alpha.1`, `-rc.3`), NPM **requires** an explicit `--tag` flag
- This prevents prerelease versions from accidentally overwriting the `latest` tag
- Without `--tag`, NPM rejects the publish to protect production users

**Our Workflow Before:**
```yaml
- name: Publish to NPM
  run: |
    cd npm-package
    npm publish --access public  # ❌ Missing --tag for beta versions
```

**What Happened:**
1. Tag `v1.3.0-beta.2` triggered workflow
2. Version extracted: `1.3.0-beta.2`
3. `package.json` updated with beta version
4. `npm publish` called **without** `--tag` flag
5. NPM detected prerelease suffix → **rejected publish**

---

## ✅ The Solution

### Workflow Changes

Added **automatic prerelease detection** with conditional publishing:

```yaml
- name: Extract version from tag
  id: get_version
  run: |
    VERSION=${GITHUB_REF#refs/tags/v}
    echo "VERSION=$VERSION" >> $GITHUB_OUTPUT
    
    # Detect prerelease (contains -, e.g., beta, alpha, rc)
    if [[ "$VERSION" == *-* ]]; then
      # Extract prerelease tag (e.g., "beta" from "1.3.0-beta.2")
      TAG=$(echo "$VERSION" | sed -E 's/.*-([a-z]+).*/\1/')
      echo "TAG=$TAG" >> $GITHUB_OUTPUT
      echo "IS_PRERELEASE=true" >> $GITHUB_OUTPUT
      echo "📦 Detected prerelease version: $VERSION (tag: $TAG)"
    else
      echo "TAG=latest" >> $GITHUB_OUTPUT
      echo "IS_PRERELEASE=false" >> $GITHUB_OUTPUT
      echo "📦 Detected stable version: $VERSION (tag: latest)"
    fi

- name: Publish to NPM (Prerelease)
  if: steps.get_version.outputs.IS_PRERELEASE == 'true'
  run: |
    cd npm-package
    npm publish --access public --tag ${{ steps.get_version.outputs.TAG }}
    echo "✅ Published as prerelease with tag: ${{ steps.get_version.outputs.TAG }}"

- name: Publish to NPM (Stable)
  if: steps.get_version.outputs.IS_PRERELEASE == 'false'
  run: |
    cd npm-package
    npm publish --access public
    echo "✅ Published as stable release with tag: latest"
```

### How It Works

**Prerelease Detection Logic:**

```bash
VERSION="1.3.0-beta.2"

# Check if version contains hyphen
if [[ "$VERSION" == *-* ]]; then
  # Extract tag: beta, alpha, rc, etc.
  TAG=$(echo "$VERSION" | sed -E 's/.*-([a-z]+).*/\1/')
  # Result: TAG="beta"
fi
```

**Examples:**

| Git Tag | Version Extracted | Detected Tag | NPM Command |
|---------|------------------|--------------|-------------|
| `v1.3.0-beta.2` | `1.3.0-beta.2` | `beta` | `npm publish --tag beta` |
| `v1.3.0-alpha.1` | `1.3.0-alpha.1` | `alpha` | `npm publish --tag alpha` |
| `v1.3.0-rc.1` | `1.3.0-rc.1` | `rc` | `npm publish --tag rc` |
| `v1.3.0` | `1.3.0` | `latest` | `npm publish` (no --tag needed) |

---

## 📚 Understanding NPM Tags

### What are NPM Tags?

NPM tags are **aliases** for specific versions. They allow users to install different release channels:

```bash
npm install gohome            # Installs "latest" tag (stable)
npm install gohome@beta       # Installs "beta" tag (prerelease)
npm install gohome@1.3.0-beta.2  # Installs exact version
```

### Default Behavior

- **Without `--tag`**: NPM publishes to `latest` tag
- **With `--tag beta`**: NPM publishes to `beta` tag (doesn't affect `latest`)

### Why This Matters

**Bad Scenario (Before Fix):**

If we published `1.3.0-beta.2` without `--tag`:
1. Users run `npm install gohome`
2. NPM gives them `1.3.0-beta.2` (beta version!)
3. Production users get unstable code 💥

**Good Scenario (After Fix):**

With `--tag beta`:
1. Users run `npm install gohome` → Get `1.2.0` (latest stable)
2. Beta testers run `npm install gohome@beta` → Get `1.3.0-beta.2`
3. Everyone is happy! ✅

---

## 🔍 About OIDC Publishing

### Configuration Status

**OIDC is correctly configured:**

```yaml
permissions:
  contents: read
  id-token: write  # ✅ Enables OIDC trusted publishing
```

**Setup Node:**
```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: '22'
    registry-url: 'https://registry.npmjs.org'  # ✅ NPM registry
```

### Why OIDC is Better

**Traditional Token Auth:**
- Requires storing `NPM_TOKEN` secret
- Token can be stolen/leaked
- Manual token rotation needed
- Token appears in logs

**OIDC Trusted Publishing:**
- ✅ No secrets to manage
- ✅ Automatic JWT token generation
- ✅ Short-lived tokens (expires in minutes)
- ✅ Cryptographically verified by NPM
- ✅ Auditablе provenance

### The "always-auth" Warning

```
npm warn Unknown user config "always-auth". This will stop working in the next major version of npm.
```

**What is this?**
- Legacy `.npmrc` config from older NPM versions
- Not needed with OIDC trusted publishing
- Can be safely ignored (workflow doesn't use it)

**Action Taken:**
- Verified no `.npmrc` file exists in repo
- Warning comes from default setup-node behavior
- Will disappear in next NPM major version

---

## 🧪 Testing the Fix

### Manual Test (Before Push)

```bash
# Test the version detection logic locally
VERSION="1.3.0-beta.2"

if [[ "$VERSION" == *-* ]]; then
  TAG=$(echo "$VERSION" | sed -E 's/.*-([a-z]+).*/\1/')
  echo "Prerelease detected: tag=$TAG"
else
  echo "Stable release: tag=latest"
fi

# Output: Prerelease detected: tag=beta
```

### Test Matrix

| Version | Expected Tag | Actual Tag | Status |
|---------|--------------|------------|--------|
| `1.3.0-beta.2` | `beta` | `beta` | ✅ |
| `1.3.0-alpha.1` | `alpha` | `alpha` | ✅ |
| `1.3.0-rc.1` | `rc` | `rc` | ✅ |
| `1.3.0` | `latest` | `latest` | ✅ |
| `2.0.0-beta.10` | `beta` | `beta` | ✅ |

### Verify After Publish

```bash
# Check NPM registry
npm view @ngockhoi96/gohome dist-tags

# Expected output:
# {
#   latest: '1.2.0',
#   beta: '1.3.0-beta.2'
# }

# Test installation
npm install @ngockhoi96/gohome       # Gets 1.2.0
npm install @ngockhoi96/gohome@beta  # Gets 1.3.0-beta.2
```

---

## 📝 Commit This Fix

### Changes Made

1. **Workflow file:** `.github/workflows/release.yml`
   - Added prerelease detection logic
   - Split publish step into conditional prerelease/stable
   - Added logging for visibility

2. **Documentation:** `.github/NPM_PRERELEASE_FIX.md` (this file)
   - Complete explanation of issue
   - Solution implementation
   - Testing guidance

### Commit Message

```bash
git add .github/workflows/release.yml .github/NPM_PRERELEASE_FIX.md
git commit -m "🐛 fix(ci): add NPM prerelease tag detection for beta/alpha/rc versions

- Auto-detect prerelease versions (contains hyphen)
- Extract prerelease tag (beta, alpha, rc) from version string
- Conditionally publish with --tag for prereleases
- Publish without --tag for stable releases
- Prevent prerelease versions from overwriting 'latest' tag

Fixes NPM publish error: 'You must specify a tag using --tag when publishing a prerelease version'

See .github/NPM_PRERELEASE_FIX.md for detailed explanation"
```

---

## 🚀 Next Steps

1. **Commit the fix:**
   ```bash
   git add -A
   git commit -m "🐛 fix(ci): add NPM prerelease tag detection"
   ```

2. **Amend or create new commit for v1.3.0-beta.2:**
   ```bash
   # Option A: Amend existing release commit
   git commit --amend --no-edit
   git tag -f v1.3.0-beta.2
   
   # Option B: Create new commit and retag
   git commit -m "fix: npm prerelease publishing"
   git tag -f v1.3.0-beta.2
   ```

3. **Push with force (to update tag):**
   ```bash
   git push origin main --force-with-lease
   git push origin v1.3.0-beta.2 --force
   ```

4. **Verify workflow:**
   - GitHub Actions will trigger on tag push
   - Check workflow logs for "📦 Detected prerelease version"
   - Verify NPM publish succeeds with `--tag beta`

---

## 📖 References

- [NPM Docs: dist-tags](https://docs.npmjs.com/cli/v10/commands/npm-dist-tag)
- [NPM Docs: Publishing packages](https://docs.npmjs.com/packages-and-modules/contributing-packages-to-the-registry)
- [GitHub Actions: OIDC Token](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [NPM OIDC Trusted Publishing](https://github.blog/2023-04-19-introducing-npm-package-provenance/)

---

## ✅ Summary

**Problem:** NPM rejected prerelease publish without `--tag` flag  
**Solution:** Auto-detect prerelease + conditional --tag based on version  
**Result:** Supports both stable (latest) and prerelease (beta/alpha/rc) publishing  
**OIDC:** Already configured correctly, no changes needed  
**Status:** Ready to push and test with v1.3.0-beta.2
