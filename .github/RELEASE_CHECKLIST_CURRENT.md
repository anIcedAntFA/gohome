# ✅ Release Checklist for v1.3.0-beta.1

## Phase 1: Pre-Release ✅ COMPLETED
- [x] Update CHANGELOG.md with v1.3.0-beta.1
- [x] Merge PR #18 into main
- [x] Create and push git tag v1.3.0-beta.1

## Phase 2: Release Creation 🔄 IN PROGRESS
- [ ] **Step 7:** Monitor GitHub Actions workflow
  - Go to: https://github.com/anIcedAntFA/gohome/actions
  - Wait for workflow to complete (~5-10 minutes)
  - Verify all jobs pass (build, test, release)
  
- [ ] **Step 8:** Edit GitHub Release
  - Go to: https://github.com/anIcedAntFA/gohome/releases/tag/v1.3.0-beta.1
  - Click "Edit release"
  - Check ✅ "Set as a pre-release" checkbox
  - Copy content from: `.github/RELEASE_NOTES_v1.3.0-beta.1.md`
  - Paste into release description
  - Click "Update release"

## Phase 3: Verification ⏳ PENDING
- [ ] **Step 9:** Test the release
  - Download binary for your platform
  - Test installation: `go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.3.0-beta.1`
  - Verify version: `gohome version`
  - Test basic functionality: `gohome report`
  - Test new features:
    - [ ] `gohome config init`
    - [ ] `gohome config show`
    - [ ] `gohome completion bash`
    - [ ] Environment variables
  
- [ ] **Step 10:** Create beta testing discussion (optional)
  ```bash
  # If you have gh CLI installed:
  gh discussion create \
    --category "Beta Testing" \
    --title "🧪 v1.3.0-beta.1 Beta Testing Feedback" \
    --body "Please share your experience testing v1.3.0-beta.1"
  
  # Or create manually:
  # https://github.com/anIcedAntFA/gohome/discussions/new?category=general
  ```

## Phase 4: Beta Testing Period (1-2 weeks)
- [ ] Monitor for bug reports
- [ ] Gather user feedback
- [ ] Fix critical issues if found
- [ ] Update documentation based on feedback

## Phase 5: Stable Release (After Beta)
- [ ] Address all beta feedback
- [ ] Update CHANGELOG (beta.1 → stable)
- [ ] Create tag v1.3.0
- [ ] Uncheck "pre-release" on GitHub
- [ ] Update package managers (Homebrew, npm, etc.)
- [ ] Announce stable release

---

## 🔗 Quick Links

- **GitHub Actions:** https://github.com/anIcedAntFA/gohome/actions
- **Releases:** https://github.com/anIcedAntFA/gohome/releases
- **Release Notes:** `.github/RELEASE_NOTES_v1.3.0-beta.1.md`
- **Issues:** https://github.com/anIcedAntFA/gohome/issues

---

## 📝 Next Immediate Steps

1. Open GitHub Actions and wait for workflow to complete
2. Once done, go to Releases page
3. Edit the v1.3.0-beta.1 release
4. Mark as pre-release and add release notes
5. Test the installation yourself
6. Share with community for testing

---

## 🎯 Success Criteria for Beta

- [ ] No critical bugs reported
- [ ] Migration from v1.2 works smoothly
- [ ] All new features work as expected
- [ ] Positive community feedback
- [ ] Documentation is clear and helpful

After meeting these criteria → Create v1.3.0 stable release!
