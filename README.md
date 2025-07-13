<div align="center">
  <img src="https://go.dev/blog/gopher/header.jpg" alt="Go Gopher" width="100%">
</div>

# GopherTube

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/dl/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)
[![Last Commit](https://img.shields.io/github/last-commit/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/commits/main)
[![Contributors](https://img.shields.io/github/contributors/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/graphs/contributors)
[![Code Size](https://img.shields.io/github/languages/code-size/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)
[![Issues](https://img.shields.io/github/issues/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/issues)
[![PRs](https://img.shields.io/github/issues-pr/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/pulls)
[![Stars](https://img.shields.io/github/stars/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/stargazers)

A modern terminal user interface for searching and watching YouTube videos using yt-dlp and mpv.

[![GopherTube CLI](https://img.shields.io/badge/GopherTube-CLI-00ADD8?style=for-the-badge&logo=terminal&logoColor=white)](https://github.com/KrishnaSSH/GopherTube)

---

## Overview

GopherTube is a terminal-based YouTube client that uses yt-dlp for searching and mpv for playback. Built with Go and Bubble Tea, it provides a clean, responsive experience for browsing and watching YouTube content directly from your terminal.

## Key Features

- **YouTube search** using yt-dlp
- **Video playback** with mpv media player
- **Terminal UI** with clean design
- **Keyboard navigation** for efficient browsing
- **Loading indicators** with spinning animations
- **YAML configuration** for customization
- **Cross-platform support** for Linux and macOS

## Recent Updates

- **Fixed video playback** - Corrected mpv command-line options
- **Added loading spinners** - Visual feedback when selecting videos
- **Improved navigation** - Better scrolling and selection behavior
- **Cleaned up code** - Removed unnecessary comments and optimizations

---

## Prerequisites

### Required Dependencies

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/dl/) **Go 1.21+** - [Download](https://golang.org/dl/)

[![mpv](https://img.shields.io/badge/mpv-Media%20Player-green?style=flat-square)](https://mpv.io/) **mpv** - Media player
```bash
# Ubuntu/Debian
sudo apt install mpv

# macOS
brew install mpv

# Arch Linux
sudo pacman -S mpv
```

[![yt-dlp](https://img.shields.io/badge/yt--dlp-YouTube%20Downloader-red?style=flat-square)](https://github.com/yt-dlp/yt-dlp) **yt-dlp** - For search and playback
```bash
pip install yt-dlp
```

### Quick Dependency Installation

```bash
# Ubuntu/Debian
sudo apt install mpv
pip install yt-dlp

# macOS
brew install mpv yt-dlp

# Arch Linux
sudo pacman -S mpv yt-dlp
```

---

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/KrishnaSSH/GopherTube.git
cd GopherTube

# Build the application
go build -o gophertube

# Run the application
./gophertube
```

### Using Make

```bash
# Build
make build

# Install
make install
```

---

## Usage

### Basic Usage

1. **Start the application:**
   ```bash
   ./gophertube
   ```

2. **Search for videos:**
   - Type your search query
   - Press Enter to search

3. **Navigate results:**
   - Use ↑/↓ arrows to move through videos
   - Press Enter to play a video
   - Press Tab to load more videos
   - Press Esc to go back to search

4. **Play videos:**
   - Select a video and press Enter
   - The video will open in mpv
   - A spinner shows while the video loads

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Search / Play selected video |
| `↑/↓` | Navigate video list |
| `Tab` | Load more videos |
| `g` | Go to first video |
| `G` | Go to last video |
| `Esc` | Go back / Quit |

---

## Configuration

Create a configuration file at `~/.config/gophertube/gophertube.yaml`:

```yaml
# Path to mpv executable
mpv_path: "mpv"

# Path to yt-dlp executable
ytdlp_path: "yt-dlp"

# Default video quality
video_quality: "best[height<=1080]/best"

# Download directory for videos
download_path: "~/Videos/gophertube"

# Number of search results to fetch
search_limit: 8
```

### Configuration Options

| Setting | Description | Default |
|---------|-------------|---------|
| `mpv_path` | Path to mpv executable | `mpv` |
| `video_quality` | Preferred video quality | `best[height<=1080]/best` |
| `ytdlp_path` | Path to yt-dlp | `yt-dlp` |
| `download_path` | Download directory | `~/Videos/gophertube` |
| `search_limit` | Maximum results per search | `8` |

---

## Project Structure

```
GopherTube/
├── main.go                    # Application entry point
├── Makefile                   # Build and installation scripts
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── LICENSE                    # GPL v3 license
├── .gitignore                 # Git ignore patterns
│
├── config/
│   └── gophertube.yaml.example  # Configuration template
│
├── man/
│   └── gophertube.1            # Manual page
│
└── internal/                   # Core application code
    ├── app/
    │   └── app.go              # Main application logic
    │
    ├── components/              # UI components
    │   ├── search.go            # Search input component
    │   ├── video_list.go        # Video list display
    │   ├── status_bar.go        # Status bar component
    │   └── types.go             # Component types
    │
    ├── services/                # Business logic
    │   ├── config.go            # Configuration management
    │   └── mpv.go               # Media player service
    │
    ├── types/                   # Shared data types
    │   └── types.go             # Core data structures
    │
    ├── interfaces/              # Interface definitions
    │   └── interfaces.go        # Service interfaces
    │
    ├── constants/               # Application constants
    │   └── constants.go         # Settings and constants
    │
    ├── utils/                   # Utility functions
    │   └── utils.go             # Helper functions
    │
    └── errors/                  # Error handling
        └── errors.go            # Custom error types
```

---

## Development

### Building

```bash
# Development build
go build -o gophertube

# Release build
make build
```

### Installing

```bash
# Install binary and man page
make install

# Install binary only
make install-binary

# Install man page only
make install-man
```

### Development Tools

This project uses modern Go development practices:

[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-TUI%20Framework-blue?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-Styling-pink?style=flat-square)](https://github.com/charmbracelet/lipgloss)
[![Cobra](https://img.shields.io/badge/Cobra-CLI%20Commands-green?style=flat-square)](https://github.com/spf13/cobra)
[![Viper](https://img.shields.io/badge/Viper-Config%20Management-purple?style=flat-square)](https://github.com/spf13/viper)

---

## License

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)

This project is licensed under the [GNU General Public License v3.0](https://github.com/KrishnaSSH/GopherTube/blob/main/LICENSE).

---

## Contributing

We welcome contributions! Please see our contributing guidelines and feel free to submit issues and pull requests.

[![Star](https://img.shields.io/github/stars/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube)
[![Fork](https://img.shields.io/github/forks/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube)
[![Watch](https://img.shields.io/github/watchers/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube) 