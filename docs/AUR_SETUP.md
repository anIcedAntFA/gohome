# AUR (Arch User Repository) Setup Guide

This document guides you through setting up AUR publishing for gohome.

## Prerequisites

- An AUR account at https://aur.archlinux.org/register
- SSH key pair for authentication
- Package name `gohome` must be available (not already taken)

## Step-by-Step Setup

### 1. Create AUR Account

1. Visit https://aur.archlinux.org/register
2. Fill in the registration form
3. Verify your email
4. Log in to your account

### 2. Generate SSH Key for AUR

Create a dedicated SSH key for AUR (don't reuse existing keys for security):

```bash
# Generate a new SSH key
ssh-keygen -t ed25519 -C "your-email@example.com" -f ~/.ssh/aur

# When prompted:
# - Enter passphrase: leave EMPTY (required for CI/CD automation)
# - Confirm passphrase: leave EMPTY
```

**Important:** The key must NOT have a passphrase for GitHub Actions to use it.

### 3. Add SSH Public Key to AUR Profile

```bash
# Display your public key
cat ~/.ssh/aur.pub
```

Copy the output and:

1. Go to https://aur.archlinux.org/account/
2. Click "My Account"
3. Paste the public key into "SSH Public Key" field
4. Click "Update"

### 4. Test SSH Connection

```bash
# Configure SSH for AUR
cat >> ~/.ssh/config << 'EOF'
Host aur.archlinux.org
  IdentityFile ~/.ssh/aur
  User aur
EOF

# Test connection
ssh aur@aur.archlinux.org
# Expected output: "Hi <username>! You've successfully authenticated..."
```

### 5. Add Private Key to GitHub Secrets

```bash
# Display private key (be careful!)
cat ~/.ssh/aur
```

Copy the **entire content** including `-----BEGIN OPENSSH PRIVATE KEY-----` and `-----END OPENSSH PRIVATE KEY-----`, then:

1. Go to https://github.com/anIcedAntFA/gohome/settings/secrets/actions
2. Click "New repository secret"
3. Name: `AUR_SSH_PRIVATE_KEY`
4. Value: Paste the private key content
5. Click "Add secret"

### 6. Register Package Name on AUR

Before publishing, you need to claim the package name:

```bash
# Clone the empty repository (will show warning - this is expected)
git clone ssh://aur@aur.archlinux.org/gohome.git aur-gohome
cd aur-gohome

# The repo is empty at this point
# Package will NOT appear on AUR website until first commit is pushed
```

**Option A: Let GoReleaser create the first commit (Recommended)**

Skip to step "Publishing Workflow". When you release a new tag, GoReleaser will automatically create and push the initial PKGBUILD.

**Option B: Manually create initial commit to verify (Optional)**

If you want to verify SSH access and see the package on AUR before releasing:

```bash
cd aur-gohome

# Create a minimal PKGBUILD (will be replaced by GoReleaser later)
cat > PKGBUILD << 'EOF'
# Maintainer: ngockhoi96 <ngockhoi96 dot dev at gmail dot com>
pkgname=gohome
pkgver=1.0.4
pkgrel=1
pkgdesc="A fast, configurable Git standup & activity reporting CLI"
arch=('x86_64' 'aarch64')
url="https://github.com/anIcedAntFA/gohome"
license=('MIT')
makedepends=('go' 'git')
source=("${pkgname}-${pkgver}.tar.gz::${url}/archive/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
  cd "${pkgname}-${pkgver}"
  go build -o gohome ./cmd/gohome
}

package() {
  cd "${pkgname}-${pkgver}"
  install -Dm755 gohome "${pkgdir}/usr/bin/gohome"
}
EOF

# Generate .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Commit and push
git add PKGBUILD .SRCINFO
git commit -m "Initial commit: gohome v1.0.4"
git push origin master

# Now check https://aur.archlinux.org/packages/gohome
# Package should appear within a few minutes
```

**Note:** If you manually push an initial commit, GoReleaser will overwrite it on the next release with the proper generated PKGBUILD.

If you get an error that the package already exists, you'll need to:
- Use a different name (e.g., `gohome-git`)
- Or submit an orphan/deletion request if it's abandoned

## Verification

After setup, verify everything works:

### Local Test (without publishing)

```bash
# Test GoReleaser AUR configuration
goreleaser release --snapshot --skip=publish --clean

# Check generated PKGBUILD
ls -la dist/gohome*.pkg.tar.zst
cat dist/aur_sources/gohome/PKGBUILD
```

### Test Installation on Arch Linux

If you have an Arch Linux system or VM:

```bash
cd dist/aur_sources/gohome

# Build and install locally
makepkg -si

# Test the installed binary
gohome --version
```

## Publishing Workflow

Once setup is complete, the release workflow will automatically:

1. Trigger when you push a tag (e.g., `v1.0.5`)
2. Build binaries with GoReleaser
3. Generate `PKGBUILD` and `.SRCINFO`
4. Commit and push to AUR repository
5. AUR users can then install with `yay -S gohome` or `paru -S gohome`

## Troubleshooting

### "Permission denied (publickey)"

- Verify SSH key is added to AUR profile
- Check `~/.ssh/config` has correct `IdentityFile`
- Test: `ssh -vvv aur@aur.archlinux.org`

### "Package name already exists"

- Check https://aur.archlinux.org/packages/gohome
- If abandoned, submit orphan request
- Otherwise, use alternative name like `gohome-git`

### "ERROR: One or more files did not pass the validity check"

- Verify checksums in `.SRCINFO`
- Run `updpkgsums` to regenerate checksums
- Make sure source tarball URL is accessible

### GoReleaser AUR publish fails

- Check `AUR_SSH_PRIVATE_KEY` is set correctly in GitHub Secrets
- Verify no passphrase on private key
- Check goreleaser output for specific error messages

## Resources

- [AUR Submission Guidelines](https://wiki.archlinux.org/title/AUR_submission_guidelines)
- [PKGBUILD Documentation](https://wiki.archlinux.org/title/PKGBUILD)
- [GoReleaser AUR Documentation](https://goreleaser.com/customization/aur_sources/)
- [AUR Helpers](https://wiki.archlinux.org/title/AUR_helpers) (yay, paru, etc.)

## Maintenance

After initial setup:

- AUR package updates automatically on new releases
- Monitor https://aur.archlinux.org/packages/gohome for user feedback
- Respond to comments and flag out-of-date notifications
- Keep dependencies up to date in `.goreleaser.yml`
