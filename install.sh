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

# Install Go with official binary
install_go() {
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "Go is already installed: $GO_VERSION"
        
        # Check if Go version is sufficient
        GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
        GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
        
        if [ "$GO_MAJOR" -gt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -ge 21 ]); then
            print_success "Go version is sufficient"
            return 0
        else
            print_warning "Go version is too old. Installing newer version..."
        fi
    fi

    print_status "Installing Go from official binary..."
    
    # Determine latest Go version
    GO_VERSION=$(curl -s https://golang.org/dl/ | grep -o 'go[0-9]\+\.[0-9]\+' | head -1 2>/dev/null || echo "go1.21.0")
    if [ -z "$GO_VERSION" ] || [ "$GO_VERSION" = "go1.21.0" ]; then
        GO_VERSION="go1.21.0"  # Fallback version
    fi
    
    print_status "Using Go version: $GO_VERSION"
    
    # Download URL
    GO_ARCH=""
    case $ARCH in
        "x86_64"|"amd64") GO_ARCH="amd64" ;;
        "i386"|"i686") GO_ARCH="386" ;;
        "armv7l"|"arm") GO_ARCH="armv6l" ;;
        "aarch64"|"arm64") GO_ARCH="arm64" ;;
        "ppc64le") GO_ARCH="ppc64le" ;;
        "s390x") GO_ARCH="s390x" ;;
        *) GO_ARCH="amd64" ;;  # Default fallback
    esac
    
    GO_OS=""
    case $OS in
        "linux") GO_OS="linux" ;;
        "macos") GO_OS="darwin" ;;
        "freebsd") GO_OS="freebsd" ;;
        "openbsd") GO_OS="openbsd" ;;
        "netbsd") GO_OS="netbsd" ;;
        *) GO_OS="linux" ;;  # Default fallback
    esac
    
    GO_URL="https://golang.org/dl/${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz"
    
    print_status "Downloading Go from: $GO_URL"
    
    # Download and install
    cd /tmp
    if curl -L "$GO_URL" -o go.tar.gz; then
        if run_with_sudo tar -C /usr/local -xzf go.tar.gz; then
            rm -f go.tar.gz
            
            # Add to PATH if not already there
            if ! grep -q "/usr/local/go/bin" ~/.bashrc 2>/dev/null; then
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            fi
            export PATH=$PATH:/usr/local/go/bin
            
            print_success "Go installed successfully"
        else
            print_error "Failed to extract Go archive"
            rm -f go.tar.gz
            return 1
        fi
    else
        print_error "Failed to download Go from $GO_URL"
        return 1
    fi
}

# Install mpv from source
install_mpv() {
    if command_exists mpv; then
        print_status "mpv is already installed"
        return 0
    fi

    print_status "Installing mpv from source..."
    
    # Check for build dependencies
    if ! command_exists git; then
        print_error "git is required to build mpv"
        return 1
    fi
    
    # Install build dependencies
    print_status "Installing build dependencies..."
    if [ "$PKG_MANAGER" != "none" ]; then
        case $PKG_MANAGER in
            "apt")
                run_with_sudo apt-get update
                run_with_sudo apt-get install -y build-essential cmake pkg-config libssl-dev libffi-dev python3-dev
                ;;
            "pacman")
                run_with_sudo pacman -S --noconfirm base-devel cmake pkg-config openssl libffi python
                ;;
            "dnf"|"yum")
                run_with_sudo dnf install -y gcc gcc-c++ make cmake pkgconfig openssl-devel libffi-devel python3-devel
                ;;
            "brew")
                brew install cmake pkg-config openssl libffi
                ;;
        esac
    fi
    
    # Clone and build mpv
    cd /tmp
    if git clone https://github.com/mpv-player/mpv.git mpv-build; then
        cd mpv-build
        ./bootstrap.py
        ./waf configure
        ./waf build
        run_with_sudo ./waf install
        print_success "mpv installed from source"
    else
        print_error "Failed to build mpv from source"
        return 1
    fi
}

# Install yt-dlp from official source
install_ytdlp() {
    if command_exists yt-dlp; then
        print_status "yt-dlp is already installed"
        return 0
    fi

    print_status "Installing yt-dlp from official source..."
    
    # Download yt-dlp directly from GitHub
    YTDLP_URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
    
    print_status "Downloading yt-dlp from: $YTDLP_URL"
    
    # Download and install
    cd /tmp
    if curl -L "$YTDLP_URL" -o yt-dlp; then
        chmod +x yt-dlp
        run_with_sudo cp yt-dlp /usr/local/bin/
        rm -f yt-dlp
        print_success "yt-dlp installed successfully"
    else
        print_error "Failed to download yt-dlp"
        return 1
    fi
}

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
    
    print_status "Building GopherTube..."
    
    # Download dependencies
    go mod download || print_warning "Failed to download some dependencies"
    
    # Build the application
    go build -o gophertube main.go || {
        print_error "Failed to build GopherTube"
        return 1
    }
    
    if [ -f "gophertube" ]; then
        print_success "GopherTube built successfully!"
        chmod +x gophertube
    else
        print_error "Build completed but binary not found"
        return 1
    fi
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
            run_with_sudo mandb 2>/dev/null || print_warning "Failed to update man database"
        fi
        
        # Create config directory and copy example
        mkdir -p "$CONFIG_DIR"
        if [ -f "config/gophertube.yaml.example" ]; then
            cp config/gophertube.yaml.example "$CONFIG_DIR/gophertube.yaml"
            print_status "Configuration file created at $CONFIG_DIR/gophertube.yaml"
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
    
    if command_exists yt-dlp; then
        print_success "yt-dlp is available"
    else
        print_warning "yt-dlp not found - search may not work"
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
        exit 1
    fi
    
    if ! command_exists git; then
        print_error "git is required but not installed"
        exit 1
    fi
    
    print_status "Starting GopherTube installation..."
    
    detect_system
    install_go
    install_mpv
    install_ytdlp
    build_gophertube
    install_system_wide
    test_installation
    
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