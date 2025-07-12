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

A modern terminal user interface for searching and watching YouTube videos with blazing-fast performance using Invidious API and mpv.

[![GopherTube CLI](https://img.shields.io/badge/GopherTube-CLI-00ADD8?style=for-the-badge&logo=terminal&logoColor=white)](https://github.com/KrishnaSSH/GopherTube)

---

## Overview

GopherTube is a high-performance terminal-based YouTube client that combines the speed of Invidious search with the reliability of yt-dlp for video playback. Built with Go and Bubble Tea, it provides a smooth, responsive experience for browsing and watching YouTube content directly from your terminal.

## Key Features

- **Lightning-fast search** using Invidious API with curl and jq
- **Seamless video playback** with mpv media player
- **Responsive terminal UI** with Swiss design principles
- **Keyboard-driven navigation** for efficient browsing
- **Smart caching** for improved performance
- **YAML configuration** for easy customization
- **Cross-platform support** for Linux and macOS

## Performance Highlights

- **Sub-second search results** via Invidious API
- **Optimized video loading** with intelligent format selection
- **Memory-efficient caching** system
- **Timeout management** for reliable operation
- **Background processing** for smooth UI experience

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

[![curl](https://img.shields.io/badge/curl-HTTP%20Client-blue?style=flat-square)](https://curl.se/) **curl** - HTTP client (usually pre-installed)

[![jq](https://img.shields.io/badge/jq-JSON%20Processor-orange?style=flat-square)](https://stedolan.github.io/jq/) **jq** - JSON processor
```bash
# Ubuntu/Debian
sudo apt install jq

# macOS
brew install jq

# Arch Linux
sudo pacman -S jq
```

### Optional Dependencies

[![yt-dlp](https://img.shields.io/badge/yt--dlp-YouTube%20Downloader-red?style=flat-square)](https://github.com/yt-dlp/yt-dlp) **yt-dlp** - For direct video downloads
```bash
pip install yt-dlp
```

### Quick Dependency Installation

```bash
# Ubuntu/Debian
sudo apt install mpv jq
pip install yt-dlp

# macOS
brew install mpv jq yt-dlp

# Arch Linux
sudo pacman -S mpv jq yt-dlp
```

---

## Installation

[![Install](https://img.shields.io/badge/Install-Ready-brightgreen?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)

### Quick Install (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/KrishnaSSH/GopherTube.git
   cd GopherTube
   ```

2. **Install system-wide:**
   ```bash
   make install
   ```

3. **Run from anywhere:**
   ```bash
   gophertube
   ```

### Manual Installation

1. **Clone and build:**
   ```bash
   git clone https://github.com/KrishnaSSH/GopherTube.git
   cd GopherTube
   go build -o gophertube
   ```

2. **Install binary (optional):**
   ```bash
   sudo cp gophertube /usr/local/bin/
   ```

3. **Install man page (optional):**
   ```bash
   make install-man
   ```

4. **Run the application:**
   ```bash
   ./gophertube
   ```

---

## Usage

[![Usage](https://img.shields.io/badge/Usage-Guide-blue?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)

### Navigation Controls

| Key | Action |
|-----|--------|
| **Tab** | Load more videos (when viewing results) |
| **Enter** | Execute search or play selected video |
| **Esc** | Go back to search view |
| **Ctrl+C** | Quit the application |
| **↑/↓** | Navigate through videos |
| **g** | Go to first video |
| **G** | Go to last video |

### Search Interface

- Type your search query in the search input
- Press **Enter** to initiate search
- Results appear instantly in the video list
- Use **Tab** to load additional results

### Video Playback

- Select a video and press **Enter** to play
- Videos open directly in mpv player
- Full mpv controls available during playback
- Automatic quality selection for optimal performance

---

## Configuration

[![Config](https://img.shields.io/badge/Config-YAML-yellow?style=flat-square)](https://yaml.org/)

Create a configuration file at `~/.config/gophertube/gophertube.yaml`:

```yaml
# Media player settings
mpv_path: "mpv"
video_quality: "best"

# Invidious API settings
invidious_instance: "https://invidious.projectsegfau.lt"
search_timeout: 10
max_search_results: 20

# Download settings (optional)
ytdlp_path: "yt-dlp"
download_path: "~/Videos/gophertube"

# Performance settings
cache_enabled: true
cache_duration: 300
```

### Configuration Options

| Setting | Description | Default |
|---------|-------------|---------|
| `mpv_path` | Path to mpv executable | `mpv` |
| `video_quality` | Preferred video quality | `best` |
| `invidious_instance` | Invidious API endpoint | `https://invidious.projectsegfau.lt` |
| `search_timeout` | Search timeout in seconds | `10` |
| `max_search_results` | Maximum results per search | `20` |
| `ytdlp_path` | Path to yt-dlp (optional) | `yt-dlp` |
| `download_path` | Download directory | `~/Videos/gophertube` |
| `cache_enabled` | Enable search caching | `true` |
| `cache_duration` | Cache duration in seconds | `300` |

---

## Project Structure

[![Structure](https://img.shields.io/badge/Structure-Organized-lightgrey?style=flat-square)](https://github.com/KrishnaSSH/GopherTube)

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
    │   └── app.go              # Main application logic and state management
    │
    ├── components/              # UI components
    │   ├── search.go            # Search input component
    │   ├── video_list.go        # Video list display component
    │   ├── status_bar.go        # Status bar component
    │   └── types.go             # Component-specific types
    │
    ├── services/                # Business logic services
    │   ├── config.go            # Configuration management
    │   ├── invidious.go         # Invidious API integration
    │   └── mpv.go              # Media player service
    │
    ├── types/                   # Shared data types
    │   └── types.go             # Core data structures
    │
    ├── interfaces/              # Interface definitions
    │   └── interfaces.go        # Service and component interfaces
    │
    ├── constants/               # Application constants
    │   └── constants.go         # Version, paths, and settings
    │
    ├── utils/                   # Utility functions
    │   └── utils.go             # Helper functions and utilities
    │
    └── errors/                  # Error handling
        └── errors.go            # Custom error types and handling
```

### Architecture Overview

- **`main.go`** - Application entry point and initialization
- **`internal/app/`** - Core application logic and state management
- **`internal/components/`** - Bubble Tea UI components
- **`internal/services/`** - Business logic and external integrations
- **`internal/types/`** - Shared data structures and types
- **`internal/interfaces/`** - Interface definitions for loose coupling
- **`internal/constants/`** - Application-wide constants and settings
- **`internal/utils/`** - Helper functions and utilities
- **`internal/errors/`** - Custom error types and error handling

---

## Development

[![Development](https://img.shields.io/badge/Development-Active-brightgreen?style=flat-square)](https://github.com/KrishnaSSH/GopherTube)

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

This project uses modern Go development practices and frameworks:

[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-TUI%20Framework-blue?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-Styling-pink?style=flat-square)](https://github.com/charmbracelet/lipgloss)
[![Cobra](https://img.shields.io/badge/Cobra-CLI%20Commands-green?style=flat-square)](https://github.com/spf13/cobra)
[![Viper](https://img.shields.io/badge/Viper-Config%20Management-purple?style=flat-square)](https://github.com/spf13/viper)

### Code Style

- Follows Go best practices and conventions
- Uses interfaces for loose coupling
- Implements proper error handling
- Includes comprehensive documentation
- Maintains clean separation of concerns

---

## Troubleshooting

[![Support](https://img.shields.io/badge/Support-Help-orange?style=flat-square)](https://github.com/KrishnaSSH/GopherTube/issues)

### Common Issues

#### Search Problems

1. **"Search not working"**
   - Check your internet connection
   - Verify curl and jq are installed
   - Try a different Invidious instance in config

2. **"Slow search results"**
   - Check your network connection
   - Try a different Invidious instance
   - Increase search timeout in config

#### Playback Issues

1. **"mpv not found"**
   - Ensure mpv is installed and in your PATH
   - Check the `mpv_path` configuration
   - Verify mpv works from command line

2. **"Video playback issues"**
   - Ensure mpv is properly installed
   - Check your internet connection
   - Verify the video URL is accessible
   - Try different video quality settings

#### Performance Issues

1. **"High memory usage"**
   - Disable caching in config
   - Reduce max search results
   - Restart the application

2. **"Slow startup"**
   - Check network connectivity
   - Verify all dependencies are installed
   - Review configuration settings

### Getting Help

- Check the [Issues](https://github.com/KrishnaSSH/GopherTube/issues) page
- Review the [man page](man/gophertube.1)
- Verify your configuration file
- Test with minimal configuration

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