.PHONY: install build

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

# Install GopherTube (build, install binary, install man page, create config)
install:
	go mod tidy
	go build -ldflags "-X main.version=$(VERSION)" -o gophertube main.go
	sudo cp gophertube /usr/local/bin/
	sudo mkdir -p /usr/local/share/man/man1
	sudo cp man/gophertube.1 /usr/local/share/man/man1/
	sudo mandb
	mkdir -p ~/.config/gophertube
	cp config/gophertube.yaml.example ~/.config/gophertube/gophertube.yaml
	@echo "GopherTube installed successfully!"
	@echo "Configuration file created at ~/.config/gophertube/gophertube.yaml"
	@echo "Run 'gophertube --help' to get started."

build:
	go build -ldflags "-X main.version=$(VERSION)" -o gophertube main.go 

	