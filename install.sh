#!/bin/bash

# GopherTube Installer Script
# This script installs all dependencies and builds GopherTube
# Compatible with most Unix-like systems

set -e  # Exit on any error

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_URL="https://github.com/KrishnaSSH/GopherTube.git"
INSTALL_DIR="/usr/local/bin"
MAN_DIR="/usr/local/share/man/man1"
CONFIG_DIR="$HOME/.config/gophertube"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_debug() {
    echo -e "${PURPLE}[DEBUG]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if running as root
is_root() {
    [ "$(id -u)" -eq 0 ]
}

# Function to check if sudo is available
has_sudo() {
    command_exists sudo
}

# Function to run command with sudo if needed
run_with_sudo() {
    if is_root; then
        "$@"
    elif has_sudo; then
        sudo "$@"
    else
        print_error "This script requires sudo privileges or root access"
        exit 1
    fi
}

# Function to install package with multiple package managers
install_package() {
    local package="$1"
    local package_name="$2"
    
    print_status "Installing $package_name..."
    
    # Try different package managers
    if command_exists apt-get; then
        run_with_sudo apt-get update
        run_with_sudo apt-get install -y "$package"
    elif command_exists pacman; then
        run_with_sudo pacman -S --noconfirm "$package"
    elif command_exists dnf; then
        run_with_sudo dnf install -y "$package"
    elif command_exists yum; then
        run_with_sudo yum install -y "$package"
    elif command_exists zypper; then
        run_with_sudo zypper install -y "$package"
    elif command_exists brew; then
        brew install "$package"
    elif command_exists pkg; then
        run_with_sudo pkg install -y "$package"
    elif command_exists apk; then
        run_with_sudo apk add "$package"
    else
        print_error "No supported package manager found"
        print_status "Please install $package_name manually"
        return 1
    fi
    
    print_success "$package_name installed successfully"
}

# Detect OS and package manager
detect_system() {
    print_status "Detecting system..."
    
    # Detect OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            OS_NAME="$NAME"
            OS_VERSION="$VERSION_ID"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        OS_NAME="macOS"
        OS_VERSION=$(sw_vers -productVersion)
    elif [[ "$OSTYPE" == "freebsd"* ]]; then
        OS="freebsd"
        OS_NAME="FreeBSD"
        OS_VERSION=$(freebsd-version)
    elif [[ "$OSTYPE" == "openbsd"* ]]; then
        OS="openbsd"
        OS_NAME="OpenBSD"
        OS_VERSION=$(uname -r)
    elif [[ "$OSTYPE" == "netbsd"* ]]; then
        OS="netbsd"
        OS_NAME="NetBSD"
        OS_VERSION=$(uname -r)
    else
        OS="unknown"
        OS_NAME="Unknown"
        OS_VERSION="Unknown"
    fi
    
    print_status "OS: $OS_NAME $OS_VERSION"
    
    # Detect architecture
    ARCH=$(uname -m)
    print_status "Architecture: $ARCH"
    
    # Detect package manager
    if command_exists apt-get; then
        PKG_MANAGER="apt"
    elif command_exists pacman; then
        PKG_MANAGER="pacman"
    elif command_exists dnf; then
        PKG_MANAGER="dnf"
    elif command_exists yum; then
        PKG_MANAGER="yum"
    elif command_exists zypper; then
        PKG_MANAGER="zypper"
    elif command_exists brew; then
        PKG_MANAGER="brew"
    elif command_exists pkg; then
        PKG_MANAGER="pkg"
    elif command_exists apk; then
        PKG_MANAGER="apk"
    else
        PKG_MANAGER="none"
        print_warning "No package manager detected"
    fi
    
    print_status "Package Manager: $PKG_MANAGER"
}

# Install required dependencies
install_package go "Go programming language"
install_package mpv "mpv media player"
install_package fzf "fzf fuzzy finder"
install_package chafa "chafa terminal image preview"

# Clone and build GopherTube
build_gophertube() {
    print_status "Setting up GopherTube..."
    
    # Determine installation directory
    if [ -d "GopherTube" ]; then
        print_warning "GopherTube directory already exists. Updating..."
        cd GopherTube
        git pull origin main || print_warning "Failed to pull latest changes"
    else
        print_status "Cloning GopherTube repository..."
        git clone "$REPO_URL" || {
            print_error "Failed to clone repository"
            return 1
        }
        cd GopherTube
    fi
    
    # Check if Go is in PATH
    if ! command_exists go; then
        print_error "Go not found in PATH"
        print_status "Please restart your terminal or run: export PATH=\$PATH:/usr/local/go/bin"
        return 1
    fi
    
    print_status "Downloading Go dependencies..."
    go mod download || print_warning "Failed to download some dependencies"
    go mod tidy || print_warning "Failed to tidy Go modules"
    
    # Build GopherTube with dynamic version from latest tag
    git fetch --tags > /dev/null 2>&1 || true
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
    print_status "Building GopherTube with version: $LATEST_TAG"
    go build -ldflags "-X main.version=$LATEST_TAG" -o gophertube main.go
    
    print_success "GopherTube installed via make install!"
}

# Install GopherTube system-wide
install_system_wide() {
    read -p "Do you want to install GopherTube system-wide? (Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        print_status "GopherTube binary is ready in the current directory."
        print_status "Run './gophertube' to start the application."
    else
        print_status "Installing GopherTube system-wide..."
        
        # Create directories
        run_with_sudo mkdir -p "$INSTALL_DIR"
        run_with_sudo mkdir -p "$MAN_DIR"
        
        # Copy binary
        run_with_sudo cp gophertube "$INSTALL_DIR/"
        
        # Copy man page if it exists
        if [ -f "man/gophertube.1" ]; then
            run_with_sudo cp man/gophertube.1 "$MAN_DIR/"
            if command_exists mandb; then
                run_with_sudo mandb 2>/dev/null || print_warning "Failed to update man database"
            else
                print_warning "'mandb' not found. Installing without updating man database. Man page may not be available."
            fi
        fi
        
        # Create config directory and copy example
        mkdir -p "$CONFIG_DIR"
        if [ -f "config/gophertube.toml" ]; then
            cp config/gophertube.toml "$CONFIG_DIR/gophertube.toml"
            print_status "Configuration file created at $CONFIG_DIR/gophertube.toml"
        fi
        
        print_success "GopherTube installed system-wide!"
        print_status "Run 'gophertube --help' to get started."
    fi
}

# Test installation
test_installation() {
    print_status "Testing installation..."
    
    # Test binary
    if [ -f "./gophertube" ]; then
        if ./gophertube --help >/dev/null 2>&1; then
            print_success "GopherTube binary works correctly"
        else
            print_warning "GopherTube binary may have issues"
        fi
    fi
    
    # Test dependencies
    if command_exists mpv; then
        print_success "mpv is available"
    else
        print_warning "mpv not found - video playback may not work"
    fi
    
    if command_exists fzf; then
        print_success "fzf is available"
    else
        print_warning "fzf not found - search may not work"
    fi

    if command_exists chafa; then
        print_success "chafa is available"
    else
        print_warning "chafa not found - image preview may not work"
    fi
}

# Main installation process
main() {
    echo "🐹 GopherTube Installer"
    echo "========================"
    echo
    
    # Check if running as root
    if is_root; then
        print_warning "Running as root. This is not recommended."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # Check for basic requirements
    if ! command_exists curl; then
        print_error "curl is required but not installed"
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    fi
    
    if ! command_exists git; then
        print_error "git is required but not installed"
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    fi
    
    print_status "Starting GopherTube installation..."
    
    detect_system || {
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    }
    install_package go "Go programming language"
    install_package mpv "mpv media player"
    install_package fzf "fzf fuzzy finder"
    install_package chafa "chafa terminal image preview"
    build_gophertube || {
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    }
    install_system_wide || {
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    }
    test_installation || {
        echo -e "\nIf this script fails, try manual installation:"
        echo "  git clone https://github.com/KrishnaSSH/GopherTube.git"
        echo "  cd GopherTube && make install"
        echo "If you still have issues, please report them on GitHub."
        exit 1
    }
    
    echo
    print_success "Installation completed successfully!"
    echo
    print_status "You can now run GopherTube:"
    if command_exists gophertube; then
        echo "  gophertube"
    else
        echo "  ./gophertube"
    fi
    echo
    print_status "For help, run: gophertube --help"
}

# Handle script interruption
trap 'print_error "Installation interrupted"; exit 1' INT TERM

# Run main function
main "$@" 
