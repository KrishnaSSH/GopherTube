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
	./update-aur.sh

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
	
