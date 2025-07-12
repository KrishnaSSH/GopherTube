.PHONY: install

# Install GopherTube (build, install binary, install man page, create config)
install:
	go mod tidy
	go build -o gophertube main.go
	sudo cp gophertube /usr/local/bin/
	sudo mkdir -p /usr/local/share/man/man1
	sudo cp man/gophertube.1 /usr/local/share/man/man1/
	sudo mandb
	mkdir -p ~/.config/gophertube
	cp config/gophertube.yaml.example ~/.config/gophertube/gophertube.yaml
	@echo "GopherTube installed successfully!"
	@echo "Configuration file created at ~/.config/gophertube/gophertube.yaml"
	@echo "Run 'gophertube --help' to get started." 

	