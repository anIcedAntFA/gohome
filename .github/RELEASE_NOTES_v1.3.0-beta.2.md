# Release Notes: v1.3.0-beta.2

**Release Date:** January 19, 2026  
**Type:** Pre-release (Beta)  
**Previous Version:** v1.3.0-beta.1

---

## 🐛 Bug Fixes & Improvements

### Config Command Error Handling

- **Consistent Error Formatting**: All commands now show errors with ❌ prefix for better UX
- **Fixed**: `gohome config <unknown>` previously showed only help text, now displays proper error message
- **Improved**: Added `SilenceErrors`/`SilenceUsage` flags for better error control
- **Enhanced**: Args validation added to catch unknown subcommands early

**Before:**
```bash
gohome report 123  → ❌ Error: invalid argument "123"
gohome config 123  → Shows help only (inconsistent)
```

**After:**
```bash
gohome report 123  → ❌ Error: invalid argument "123"  
gohome config 123  → ❌ Error: invalid argument "123" (consistent!)
```

### Config List Output Enhancement

- ⚙️ Added emoji to 'Current Configuration' header for better visibility
- 📄 Added emoji to config file path display
- 🧹 Removed unused table options (nature/tech type) to reduce clutter

---

## 📚 Documentation Updates

- ✅ Added comprehensive **release checklist template** (`.github/RELEASE_CHECKLIST_CURRENT.md`)
- ✅ Added **release notes template** for v1.3.0-beta.1 as reference
- ✅ Updated **README** with improved examples and clarity
- ✅ Updated **CLI_GUIDE** with latest command usage patterns
- ✅ Added **beta release testing script** (`scripts/test-beta-release.sh`) for QA automation
- ✅ Refreshed **quickstart demo** with better visual presentation

---

## 🔄 What Changed Since v1.3.0-beta.1

### Commits in This Release

1. **📝 docs(docs,readme)**: Update documents and assets
2. **🐛 fix(cmd)**: Add emojis to config list output and fix error consistency
3. **📝 docs(changelog,scripts)**: Update changelog, add script test beta version

### Files Modified

- `.github/RELEASE_CHECKLIST_CURRENT.md` (new)
- `.github/RELEASE_NOTES_v1.3.0-beta.1.md` (new)
- `CHANGELOG.md`
- `README.md`
- `cmd/gohome/cmd/config.go`
- `cmd/gohome/cmd/report.go`
- `cmd/gohome/cmd/root.go`
- `docs/CLI_GUIDE.md`
- `scripts/test-beta-release.sh` (new)

---

## 🧪 Testing

### Manual Testing

Run the beta release testing script:

```bash
bash scripts/test-beta-release.sh
```

### Key Test Cases

1. ✅ Config command error handling consistency
2. ✅ Config list output with emojis
3. ✅ Report command error validation
4. ✅ Version display: `gohome version`

---

## 📦 Installation

### Via Script (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/anIcedAntFA/gohome/v1.3.0-beta.2/scripts/install.sh | bash -s -- --beta
```

### Via Go Install

```bash
go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.3.0-beta.2
```

### From Release Assets

Download the pre-built binary for your platform from the Assets section below.

---

## 🚨 Breaking Changes

**None.** This is a maintenance release with bug fixes and documentation improvements.

---

## 🐛 Known Issues

- None reported for this beta release

---

## 📝 Migration from v1.3.0-beta.1

No migration required. This is a drop-in replacement with:
- ✅ Backward compatible CLI
- ✅ Same config file format
- ✅ Identical command structure

Simply update to beta.2 using any installation method above.

---

## 🤝 Feedback

This is a **pre-release** for testing purposes. Please report any issues:

- 🐛 [Report Bugs](https://github.com/anIcedAntFA/gohome/issues/new?labels=bug)
- 💡 [Request Features](https://github.com/anIcedAntFA/gohome/issues/new?labels=enhancement)
- 💬 [Discussions](https://github.com/anIcedAntFA/gohome/discussions)

---

## 📅 Roadmap to v1.3.0 Stable

- [x] Beta 1: Core Cobra/Viper migration
- [x] **Beta 2: Bug fixes and documentation** ← **You are here**
- [ ] Beta 3: Additional testing and polish (if needed)
- [ ] RC1: Release candidate with final testing
- [ ] v1.3.0: Stable release

---

## 🙏 Credits

Thanks to all contributors and testers helping to make gohome better!

**Full Changelog**: https://github.com/anIcedAntFA/gohome/compare/v1.3.0-beta.1...v1.3.0-beta.2
