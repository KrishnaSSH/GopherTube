.PHONY: build install clean aur-update aur-push

# Build the binary
build:
	go build -o gophertube main.go

# Install to system
install: build
	sudo cp gophertube /usr/local/bin/
	sudo chmod +x /usr/local/bin/gophertube

# Clean build artifacts
clean:
	rm -f gophertube
	rm -rf aur-gophertube

# Update AUR package with current version
aur-update:
	@echo "Updating AUR package..."
	@VERSION=$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/v//'); \
	if [ -z "$$VERSION" ]; then \
		echo "No version tag found. Please create a tag first."; \
		exit 1; \
	fi; \
	echo "Current version: $$VERSION"; \
	\
	# Download source and calculate checksum
	wget -O gophertube-$$VERSION.tar.gz https://github.com/KrishnaSSH/GopherTube/archive/refs/tags/v$$VERSION.tar.gz; \
	CHECKSUM=$$(sha256sum gophertube-$$VERSION.tar.gz | cut -d' ' -f1); \
	echo "Checksum: $$CHECKSUM"; \
	\
	# Clone AUR repository
	git clone ssh://aur@aur.archlinux.org/gophertube.git aur-gophertube; \
	cd aur-gophertube; \
	\
	# Update PKGBUILD
	sed -i "s/pkgver=.*/pkgver=$$VERSION/" PKGBUILD; \
	sed -i "s/sha256sums=('.*')/sha256sums=('$$CHECKSUM')/" PKGBUILD; \
	\
	# Generate .SRCINFO
	makepkg --printsrcinfo > .SRCINFO; \
	\
	# Commit and push
	git config --local user.email "krishna.pytech@gmail.com"; \
	git config --local user.name "KrishnaSSH"; \
	git add PKGBUILD .SRCINFO; \
	git commit -m "Update to v$$VERSION"; \
	git push origin master; \
	\
	# Cleanup
	cd ..; \
	rm -rf aur-gophertube gophertube-$$VERSION.tar.gz; \
	echo "AUR package updated to v$$VERSION"

# Push new tag and update AUR
aur-push:
	@echo "Creating and pushing new tag..."
	@VERSION=$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/v//'); \
	if [ -z "$$VERSION" ]; then \
		echo "No version tag found. Please create a tag first."; \
		exit 1; \
	fi; \
	echo "Current version: $$VERSION"; \
	\
	# Push tag if not already pushed
	git push origin v$$VERSION 2>/dev/null || echo "Tag already pushed"; \
	\
	# Update AUR
	$(MAKE) aur-update

# Show current version
version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "No version tag found"

# Create new version tag
tag:
	@echo "Usage: make tag VERSION=2.2.1"
	@if [ -z "$(VERSION)" ]; then \
		echo "Please specify VERSION=2.2.1"; \
		exit 1; \
	fi; \
	git tag v$(VERSION); \
	git push origin v$(VERSION); \
	echo "Created and pushed tag v$(VERSION)"

# Full release process
release:
	@echo "Usage: make release VERSION=2.2.1"
	@if [ -z "$(VERSION)" ]; then \
		echo "Please specify VERSION=2.2.1"; \
		exit 1; \
	fi; \
	$(MAKE) tag VERSION=$(VERSION); \
	$(MAKE) aur-update; \
	echo "Release v$(VERSION) complete!"
	
