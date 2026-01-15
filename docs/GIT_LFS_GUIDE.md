# Git LFS Guide for gohome

## Overview

Starting from **v1.2.0**, gohome uses **Git LFS (Large File Storage)** to manage demo media files (GIF, PNG, MP4, WebM). This keeps the repository lightweight while allowing scalable media content for documentation.

## Why Git LFS?

**Before LFS (v1.0-v1.1):**
- Demo GIFs stored as regular Git objects
- Total: 4.3MB (2.4MB + 1.9MB)
- Every clone downloads full media history
- Adding more demos bloats repository size

**With LFS (v1.2+):**
- Media files replaced with tiny pointers (132 bytes each)
- Actual files stored on GitHub's LFS servers
- Clone downloads pointers only (~21MB .git size)
- Media fetched on-demand or during checkout with `lfs: true`
- Scalable for unlimited demo additions

## How It Works

### File Flow

```
┌─────────────────┐
│  git add        │
│  demo.gif       │
└────────┬────────┘
         │
         ▼
┌─────────────────────┐
│ Git LFS Clean Filter│  ← Convert file → pointer
│  (.gitattributes)   │
└────────┬────────────┘
         │
         ▼
┌─────────────────┐       ┌──────────────────┐
│ Git Repository  │       │ GitHub LFS Server│
│  (pointer file) │       │  (actual GIF)    │
│  132 bytes      │       │  2.4MB          │
└─────────────────┘       └──────────────────┘
         │                         │
         └─────────┬───────────────┘
                   │
                   ▼
           ┌───────────────┐
           │ git checkout  │
           │ (lfs: true)   │
           └───────┬───────┘
                   │
                   ▼
           ┌──────────────────┐
           │ Git LFS Smudge   │  ← Download actual file
           │    Filter        │
           └───────┬──────────┘
                   │
                   ▼
           ┌──────────────┐
           │ Working Dir  │
           │  demo.gif    │
           │  2.4MB       │
           └──────────────┘
```

### Tracked Patterns

From [`.gitattributes`](../.gitattributes):

```gitattributes
docs/demos/*.gif filter=lfs diff=lfs merge=lfs -text
docs/demos/*.png filter=lfs diff=lfs merge=lfs -text
docs/demos/*.mp4 filter=lfs diff=lfs merge=lfs -text
docs/demos/*.webm filter=lfs diff=lfs merge=lfs -text
```

**Explanation:**
- `filter=lfs`: Use LFS clean/smudge filters
- `diff=lfs`: Use LFS diff driver
- `merge=lfs`: Use LFS merge driver
- `-text`: Treat as binary (no line ending conversion)

## For Contributors

### Prerequisites

Install Git LFS:

```bash
# macOS (Homebrew)
brew install git-lfs

# Ubuntu/Debian
sudo apt-get install git-lfs

# Arch Linux
sudo pacman -S git-lfs

# Fedora/RedHat
sudo dnf install git-lfs

# Windows (with Git for Windows)
# LFS is included, just run:
git lfs install
```

### Cloning the Repository

**Option 1: Clone with LFS (Recommended for contributors)**

```bash
git clone https://github.com/anIcedAntFA/gohome.git
cd gohome

# LFS is auto-initialized on clone
# Media files are automatically fetched
ls -lh docs/demos/*.gif
# -rw-r--r-- 1 user user 2.4M Jan 12 10:06 config.gif
# -rw-r--r-- 1 user user 1.9M Jan 12 10:06 quickstart.gif
```

**Option 2: Clone without LFS (Lightweight)**

For developers who don't need demo media:

```bash
GIT_LFS_SKIP_SMUDGE=1 git clone https://github.com/anIcedAntFA/gohome.git
cd gohome

# LFS files are pointer files
ls -lh docs/demos/*.gif
# -rw-r--r-- 1 user user 132 Jan 15 13:26 config.gif  ← pointer
# -rw-r--r-- 1 user user 132 Jan 15 13:26 quickstart.gif  ← pointer

# Fetch later if needed:
git lfs pull
```

### Adding New Demo Files

When adding new GIFs or media to `docs/demos/`:

```bash
# 1. Create your demo GIF (using VHS, etc.)
vhs docs/demos/new-demo.tape

# 2. Git add normally - LFS handles it automatically
git add docs/demos/new-demo.gif

# 3. Verify it's tracked by LFS
git lfs ls-files
# cd80c5fd53 * docs/demos/config.gif
# 14617d61e5 * docs/demos/quickstart.gif
# abc1234567 * docs/demos/new-demo.gif  ← Your new file

# 4. Commit and push
git commit -m "📹 docs(demos): add new-demo example"
git push origin your-branch
```

**Important:** LFS files are uploaded during `git push`, not `git commit`. Large files may take longer to push.

### Checking LFS Status

```bash
# List all LFS-tracked files
git lfs ls-files

# Show LFS storage usage
git lfs status

# Show LFS pointer details
git lfs pointer --file=docs/demos/config.gif
```

### Common Issues

#### Issue 1: "This repository is over its data quota"

GitHub free accounts have **1GB LFS storage** and **1GB bandwidth/month**.

**Solutions:**
- Wait until next month for bandwidth reset
- Upgrade to GitHub Pro (50GB storage, 50GB bandwidth)
- Use `GIT_LFS_SKIP_SMUDGE=1` to skip LFS downloads

#### Issue 2: LFS files showing as pointer text

```bash
$ cat docs/demos/config.gif
version https://git-lfs.github.com/spec/v1
oid sha256:cd80c5fd53...
size 2467212
```

**Solution:** Fetch LFS files

```bash
git lfs pull
```

#### Issue 3: "Encountered X file(s) that should have been pointers"

This happens if you commit a large file before LFS tracking was set up.

**Solution:** Re-add the file after LFS setup

```bash
git rm --cached docs/demos/file.gif
git add docs/demos/file.gif
git commit -m "fix: track file.gif with LFS"
```

## For CI/CD

### GitHub Actions

Both [release.yml](../.github/workflows/release.yml) and [build.yml](../.github/workflows/build.yml) have `lfs: true` enabled:

```yaml
- name: Checkout code
  uses: actions/checkout@v4
  with:
    fetch-depth: 0
    lfs: true  # ← Fetch LFS files during checkout
```

**Why this is important:**
- Without `lfs: true`, workflows only see 132-byte pointer files
- GoReleaser needs actual GIF files for packaging
- Build validation needs actual files to test demos

### Local Testing

Test your changes before pushing:

```bash
# Clean clone to simulate CI environment
cd /tmp
git clone https://github.com/anIcedAntFA/gohome.git test-lfs
cd test-lfs

# Verify LFS files are fetched
ls -lh docs/demos/*.gif
file docs/demos/*.gif

# Run build
make build
./bin/gohome --version
```

## Migration History

| Version | Storage | Details |
|---------|---------|---------|
| v1.0-v1.1 | Regular Git | GIFs stored as Git blobs (4.3MB total) |
| v1.2.0+ | Git LFS | GIFs converted to LFS pointers (264 bytes total) |

**Migration Commits:**
- `d9c06ec`: Setup Git LFS tracking in `.gitattributes`
- `10f4331`: Migrate existing GIFs to LFS storage
- `49b503a`: Enable LFS in GitHub Actions workflows

## Best Practices

### For New Media Files

1. **Keep demos concise**: Aim for <5MB per GIF
2. **Optimize before committing**:
   ```bash
   # Using gifsicle
   gifsicle -O3 --colors 256 input.gif -o output.gif
   
   # Using VHS (built-in optimization)
   vhs demo.tape  # Already optimized
   ```
3. **Use appropriate formats**:
   - `.gif`: Short terminal demos (2-5MB)
   - `.mp4`: Longer videos (better compression)
   - `.webm`: Web-optimized videos

### For Repository Maintainers

1. **Monitor LFS bandwidth**: Check [GitHub Settings → Billing](https://github.com/settings/billing)
2. **Prune old LFS objects** (after major version bumps):
   ```bash
   git lfs prune --verify-remote
   ```
3. **Document media in PR descriptions**: Mention file sizes and purpose

## References

- [Git LFS Official Docs](https://git-lfs.github.com/)
- [GitHub LFS Billing](https://docs.github.com/en/billing/managing-billing-for-git-large-file-storage/about-billing-for-git-large-file-storage)
- [Git LFS Tutorial](https://www.atlassian.com/git/tutorials/git-lfs)
- [VHS (Demo Recording Tool)](https://github.com/charmbracelet/vhs)

## FAQ

**Q: Will cloning be slower with LFS?**  
A: Initial clone is similar. Subsequent pulls only download changed LFS files.

**Q: Can I disable LFS locally?**  
A: Yes, use `GIT_LFS_SKIP_SMUDGE=1` before clone or `git lfs install --skip-smudge`.

**Q: What happens if GitHub LFS quota is exceeded?**  
A: LFS downloads fail, but Git operations (clone, pull) still work. Pointer files are checked out instead.

**Q: How do I check my LFS bandwidth usage?**  
A: Go to [GitHub Settings → Billing → Git LFS](https://github.com/settings/billing) (for free accounts: 1GB/month).

**Q: Can I migrate back to regular Git?**  
A: Yes, but requires history rewrite (not recommended):
```bash
git lfs migrate export --include="docs/demos/*.gif" --everything
```

---

**Need help?** Open an issue at https://github.com/anIcedAntFA/gohome/issues
