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
[![Lines of Code](https://img.shields.io/tokei/lines/github/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)
[![Issues](https://img.shields.io/github/issues/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/issues)
[![PRs](https://img.shields.io/github/issues-pr/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/pulls)
[![Stars](https://img.shields.io/github/stars/KrishnaSSH/GopherTube?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube/stargazers)

A terminal user interface for searching and watching YouTube videos using yt-dlp and mpv.

[![GopherTube CLI](https://img.shields.io/badge/GopherTube-CLI-00ADD8?style=for-the-badge&logo=terminal&logoColor=white)](https://github.com/KrishnaSSH/GopherTube)

## Features

- Terminal-based user interface with Swiss design principles
- YouTube video search with real-time results
- Direct video playback with mpv
- Keyboard navigation for video browsing
- Responsive layout that adapts to terminal size
- YAML-based configuration system

## Prerequisites

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/dl/) **Go 1.21+** - [Download](https://golang.org/dl/)

[![yt-dlp](https://img.shields.io/badge/yt--dlp-Latest-orange?style=flat-square)](https://github.com/yt-dlp/yt-dlp) **yt-dlp** - YouTube downloader
```bash
pip install yt-dlp
```

[![mpv](https://img.shields.io/badge/mpv-Media%20Player-green?style=flat-square)](https://mpv.io/) **mpv** - Media player
```bash
# Ubuntu/Debian
sudo apt install mpv

# macOS
brew install mpv

# Arch Linux
sudo pacman -S mpv
```

### Dependencies Installation

Install all dependencies at once:

```bash
# Ubuntu/Debian
sudo apt install mpv
pip install yt-dlp

# macOS
brew install mpv yt-dlp

# Arch Linux
sudo pacman -S mpv yt-dlp
```

## Installation

[![Install](https://img.shields.io/badge/Install-Ready-brightgreen?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)

### Quick Install (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/KrishnaSSH/GopherTube.git
   cd GopherTube
   ```

2. Install system-wide:
   ```bash
   make install
   ```

3. Run from anywhere:
   ```bash
   gophertube
   ```

### Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/KrishnaSSH/GopherTube.git
   cd GopherTube
   ```

2. Build the application:
   ```bash
   go build -o gophertube
   ```

3. Install binary (optional):
   ```bash
   sudo cp gophertube /usr/local/bin/
   ```

4. Install man page (optional):
   ```bash
   make install-man
   ```

5. Run the application:
   ```bash
   ./gophertube
   ```

## Usage

[![Usage](https://img.shields.io/badge/Usage-Guide-blue?style=for-the-badge)](https://github.com/KrishnaSSH/GopherTube)

### Navigation

- **Tab** - Load more videos (when viewing results)
- **Enter** - Execute search or play selected video
- **Esc** - Go back to search view
- **Ctrl+C** - Quit the application

### Search View

- Type your search query in the search box
- Press **Enter** to search
- Results will appear in the video list

### Video List View

- **↑/↓** - Navigate through videos
- **Enter** - Play the selected video
- **g** - Go to first video
- **G** - Go to last video
- **Tab** - Load more videos
- **Esc** - Return to search

## Configuration

[![Config](https://img.shields.io/badge/Config-YAML-yellow?style=flat-square)](https://yaml.org/)

Create a configuration file at `~/.config/gophertube/gophertube.yaml`:

```yaml
# Path to mpv executable
mpv_path: "mpv"

# Path to yt-dlp executable
ytdlp_path: "yt-dlp"

# Default video quality
video_quality: "best"

# Download directory
download_path: "~/Videos/gophertube"
```

## Project Structure

[![Structure](https://img.shields.io/badge/Structure-Organized-lightgrey?style=flat-square)](https://github.com/KrishnaSSH/GopherTube)

```
gophertube/
├── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go         # Main application logic
│   ├── components/
│   │   ├── search.go       # Search component
│   │   ├── video_list.go   # Video list component
│   │   ├── status_bar.go   # Status bar component
│   │   └── types.go        # Shared types
│   └── services/
│       ├── config.go       # Configuration service
│       └── mpv.go          # MPV player service
```

## Development

[![Development](https://img.shields.io/badge/Development-Active-brightgreen?style=flat-square)](https://github.com/KrishnaSSH/GopherTube)

### Building

```bash
go build -o gophertube
```

### Installing

```bash
make install
```

### Code Style

This project follows Go best practices and uses:

[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-TUI%20Framework-blue?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-Styling-pink?style=flat-square)](https://github.com/charmbracelet/lipgloss)
[![Cobra](https://img.shields.io/badge/Cobra-CLI%20Commands-green?style=flat-square)](https://github.com/spf13/cobra)
[![Viper](https://img.shields.io/badge/Viper-Config%20Management-purple?style=flat-square)](https://github.com/spf13/viper)

## Troubleshooting

[![Support](https://img.shields.io/badge/Support-Help-orange?style=flat-square)](https://github.com/KrishnaSSH/GopherTube/issues)

### Common Issues

1. **"mpv not found"**
   - Ensure mpv is installed and in your PATH
   - Check the `mpv_path` configuration

2. **"yt-dlp not found"**
   - Install yt-dlp: `pip install yt-dlp`
   - Check the `ytdlp_path` configuration

3. **Video playback issues**
   - Ensure mpv is properly installed
   - Check your internet connection
   - Verify the video URL is accessible

4. **Search not working**
   - Check your internet connection
   - Verify yt-dlp is properly installed

## License

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)

This project is licensed under the [GNU GPLv3](https://github.com/KrishnaSSH/GopherTube/blob/main/LICENSE).

---

[![Star](https://img.shields.io/github/stars/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube)
[![Fork](https://img.shields.io/github/forks/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube)
[![Watch](https://img.shields.io/github/watchers/KrishnaSSH/GopherTube?style=social)](https://github.com/KrishnaSSH/GopherTube) 