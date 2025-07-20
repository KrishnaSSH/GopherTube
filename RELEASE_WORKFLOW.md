# Automated Release Workflow

This repository includes automated workflows that handle the complete release process when you push a new tag.

## 🚀 How It Works

When you push a new tag (e.g., `v2.3.0`), the following happens automatically:

1. **GitHub Release**: Creates/updates GitHub release with DEB and RPM packages
2. **Arch Repository**: Updates the main PKGBUILD for the official Arch repository
3. **AUR Package**: Updates the AUR package (requires SSH key setup)

## 📋 Setup Required

### 1. GitHub Secrets

You need to add these secrets to your GitHub repository:

1. Go to your repository → Settings → Secrets and variables → Actions
2. Add the following secret:

#### `AUR_SSH_KEY`
- Generate an SSH key pair for AUR access
- Add the **private key** as the secret value
- The public key should be added to your AUR account

### 2. AUR SSH Key Setup

```bash
# Generate SSH key pair
ssh-keygen -t rsa -b 4096 -C "your-email@example.com" -f ~/.ssh/aur_key

# Add public key to AUR account
cat ~/.ssh/aur_key.pub
# Copy this to your AUR account settings

# Add private key to GitHub secrets
cat ~/.ssh/aur_key
# Copy this to GitHub repository secrets as AUR_SSH_KEY
```

## 🏷️ How to Release

### Simple Release Process

1. **Create and push a new tag:**
   ```bash
   git tag v2.3.0
   git push origin v2.3.0
   ```

2. **That's it!** The workflows will automatically:
   - Build DEB and RPM packages
   - Create/update GitHub release
   - Update Arch repository PKGBUILD
   - Update AUR package

### Manual Steps (if needed)

If the automated workflows fail, you can manually:

1. **Update GitHub Release:**
   - Go to GitHub → Releases
   - Edit the release for your tag
   - Upload the built packages

2. **Update Arch PKGBUILD:**
   ```bash
   # Update version and checksum in PKGBUILD
   sed -i 's/pkgver=.*/pkgver=2.3.0/' PKGBUILD
   # Calculate new checksum and update
   git add PKGBUILD && git commit -m "Update to v2.3.0"
   ```

3. **Update AUR Package:**
   ```bash
   git clone ssh://aur@aur.archlinux.org/gophertube.git
   cd gophertube
   # Update PKGBUILD and .SRCINFO
   git add . && git commit -m "Update to v2.3.0"
   git push origin master
   ```

## 🔧 Workflow Files

- **`.github/workflows/release.yml`**: Main release workflow
- **`.github/workflows/aur-update.yml`**: AUR package update workflow

## 📦 What Gets Released

- **GitHub**: DEB and RPM packages for Linux
- **Arch Repository**: PKGBUILD for official Arch packages
- **AUR**: Package for Arch User Repository

## 🐛 Troubleshooting

### AUR Update Fails
- Check that `AUR_SSH_KEY` secret is set correctly
- Verify SSH key is added to your AUR account
- Check AUR repository permissions

### GitHub Release Fails
- Ensure tag doesn't already have a release
- Check repository permissions for Actions

### Arch PKGBUILD Update Fails
- Verify GitHub token permissions
- Check if PKGBUILD format is correct

## 📝 Notes

- The AUR workflow runs after the main release workflow completes
- All workflows use the same version from the git tag
- Checksums are automatically calculated and updated
- The process is fully automated once set up correctly 