#!/bin/bash

echo "Updating AUR package..."

# Get current version
VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/v//')
if [ -z "$VERSION" ]; then
    echo "No version tag found. Please create a tag first."
    exit 1
fi

echo "Current version: $VERSION"

# Download source and calculate checksum
wget -O gophertube-$VERSION.tar.gz https://github.com/KrishnaSSH/GopherTube/archive/refs/tags/v$VERSION.tar.gz
CHECKSUM=$(sha256sum gophertube-$VERSION.tar.gz | cut -d' ' -f1)
echo "Checksum: $CHECKSUM"

# Clone AUR repository
git clone ssh://aur@aur.archlinux.org/gophertube.git aur-gophertube
cd aur-gophertube

# Update PKGBUILD
sed -i "s/pkgver=.*/pkgver=$VERSION/" PKGBUILD
sed -i "s/sha256sums=('.*')/sha256sums=('$CHECKSUM')/" PKGBUILD

# Generate .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Commit and push
git config --local user.email "krishna.pytech@gmail.com"
git config --local user.name "KrishnaSSH"
git add PKGBUILD .SRCINFO
git commit -m "Update to v$VERSION"
git push origin master

# Cleanup
cd ..
rm -rf aur-gophertube gophertube-$VERSION.tar.gz
echo "AUR package updated to v$VERSION" 