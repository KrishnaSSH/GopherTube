#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

print_status() { echo -e "${GREEN}[INFO]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Install dependencies
install_deps() {
    print_status "Installing dependencies..."
    
    if command -v apt-get >/dev/null; then
        sudo apt-get update && sudo apt-get install -y go mpv fzf chafa
    elif command -v pacman >/dev/null; then
        sudo pacman -S --noconfirm go mpv fzf chafa
    elif command -v dnf >/dev/null; then
        sudo dnf install -y go mpv fzf chafa
    elif command -v yum >/dev/null; then
        sudo yum install -y go mpv fzf chafa
    elif command -v zypper >/dev/null; then
        sudo zypper install -y go mpv fzf chafa
    elif command -v brew >/dev/null; then
        brew install go mpv fzf chafa
    elif command -v pkg >/dev/null; then
        sudo pkg install -y go mpv fzf chafa
    elif command -v apk >/dev/null; then
        sudo apk add go mpv fzf chafa
    else
        print_error "Unsupported package manager. Please perform manual installation."
        exit 1
    fi
}

# Build and install
build_install() {
    print_status "Building GopherTube..."
    
    [ ! -d "GopherTube" ] && git clone https://github.com/KrishnaSSH/GopherTube.git
    cd GopherTube
    
    go mod download
    go build -o gophertube .
    
    sudo cp gophertube /usr/local/bin/
    sudo mkdir -p /usr/local/share/man/man1
    sudo cp man/gophertube.1 /usr/local/share/man/man1/ 2>/dev/null || true
    
    mkdir -p ~/.config/gophertube
    cp config/gophertube.toml ~/.config/gophertube/ 2>/dev/null || true
    
    print_status "GopherTube installed successfully!"
}

install_deps
build_install 
