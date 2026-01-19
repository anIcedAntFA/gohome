# Release Checklist: v1.3.0-beta.2

**Date:** 2026-01-19  
**Type:** Pre-release (Beta)

---

## ✅ Completed Steps

- [x] Check changes since v1.3.0-beta.1
- [x] Update CHANGELOG.md with beta.2 section
- [x] Bump version in cmd/gohome/cmd/root.go to v1.3.0-beta.2
- [x] Create release commit: `🔖 chore(release): bump version to v1.3.0-beta.2`
- [x] Create annotated git tag: `v1.3.0-beta.2`
- [x] Create release notes: `.github/RELEASE_NOTES_v1.3.0-beta.2.md`

---

## 🚀 Next Steps (To Deploy)

### 1. Push to GitHub

```bash
# Push commits
git push origin main

# Push tag
git push origin v1.3.0-beta.2
```

### 2. Create GitHub Pre-release

Go to: https://github.com/anIcedAntFA/gohome/releases/new

- **Tag:** `v1.3.0-beta.2`
- **Title:** `v1.3.0-beta.2 - Bug Fixes & Documentation`
- **Description:** Copy from `.github/RELEASE_NOTES_v1.3.0-beta.2.md`
- **Pre-release:** ✅ Check "This is a pre-release"
- **Publish**

### 3. Verify Build

- Wait for GitHub Actions to complete
- Check release assets are built for all platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64, arm64)

### 4. Test Installation

```bash
# Test script installation
curl -sSL https://raw.githubusercontent.com/anIcedAntFA/gohome/v1.3.0-beta.2/scripts/install.sh | bash -s -- --beta

# Verify version
gohome version

# Run beta testing script
bash scripts/test-beta-release.sh
```

---

## 📋 Changes Summary

### Since v1.3.0-beta.1

**Commits:**
- `266ec1a` - docs(docs,readme): update documents and assets
- `7bef388` - fix(cmd): add emojis to config list output and fix error consistency
- `1425904` - docs(changelog,scripts): update changelog, add script test beta version
- `1fc3759` - chore(release): bump version to v1.3.0-beta.2 ← **Release commit**

**Key Improvements:**
- ✅ Fixed config command error handling
- ✅ Added emojis to config list output
- ✅ Updated documentation
- ✅ Added beta testing script

---

## 🔍 Pre-Push Verification

Run these checks before pushing:

```bash
# Verify tag exists
git tag | grep v1.3.0-beta.2

# Verify commit message
git log -1 --pretty=format:"%s"

# Verify changelog updated
git diff v1.3.0-beta.1 CHANGELOG.md

# Verify version bump
grep -n "v1.3.0-beta.2" cmd/gohome/cmd/root.go

# Build locally
make build
./bin/gohome version
```

---

## 📝 Notes

- Release notes template created at: `.github/RELEASE_NOTES_v1.3.0-beta.2.md`
- This is a **maintenance release** with no breaking changes
- Safe to deploy immediately after testing
- Users can upgrade from beta.1 without migration

---

## 🎯 Success Criteria

- [ ] GitHub release published with pre-release flag
- [ ] All platform binaries built successfully
- [ ] Installation script works with `--beta` flag
- [ ] Version command shows `v1.3.0-beta.2`
- [ ] No breaking changes reported

---

## 🐛 Rollback Plan (If Needed)

If issues are discovered:

1. Do NOT delete the tag
2. Create a new beta.3 with fixes
3. Update release notes to mark beta.2 as deprecated
4. Reference beta.2 issues in beta.3 changelog
