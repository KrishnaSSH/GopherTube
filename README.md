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

A simple terminal YouTube client for searching and watching videos using [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [mpv](https://mpv.io/).

---

## Overview

GopherTube is a terminal-based YouTube client. It uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) to search YouTube and [mpv](https://mpv.io/) to play videos. The UI is built with [Go](https://golang.org/) and [Bubble Tea](https://github.com/charmbracelet/bubbletea), and is fully keyboard-driven.

## Features

- Search YouTube with [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- Play videos with [mpv](https://mpv.io/)
- Minimal terminal UI ([Bubble Tea](https://github.com/charmbracelet/bubbletea))
- Keyboard navigation (arrows, Enter, Tab, Esc, g/G)
- Spinner/loading indicator when opening videos
- [YAML](https://yaml.org/) config for paths and settings

---

## Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [mpv](https://mpv.io/) (media player)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) (YouTube downloader)

Install dependencies:

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

```bash
git clone https://github.com/KrishnaSSH/GopherTube.git
cd GopherTube
go build -o gophertube
./gophertube
```

Or use the [Makefile](Makefile):

```bash
make build   # Build the binary
make install # Install binary and man page
```

---

## Usage

- Start the app: `./gophertube`
- Type a search and press Enter
- Use ↑/↓ to move, Enter to play, Tab to load more, Esc to go back
- Spinner shows while video loads, then [mpv](https://mpv.io/) opens

### Keyboard Shortcuts

| Key      | Action                  |
|----------|-------------------------|
| Enter    | Search / Play video     |
| ↑/↓      | Navigate video list     |
| Tab      | Load more videos        |
| g        | Go to first video       |
| G        | Go to last video        |
| Esc      | Go back / Quit          |

---

## Configuration

Create `~/.config/gophertube/gophertube.yaml`:

```yaml
mpv_path: "mpv"
ytdlp_path: "yt-dlp"
video_quality: "best[height<=1080]/best"
download_path: "~/Videos/gophertube"
search_limit: 8
```

See [config/gophertube.yaml.example](config/gophertube.yaml.example) for the full template.

---

## Project Structure

```
GopherTube/
├── main.go
├── Makefile
├── go.mod
├── go.sum
├── LICENSE
├── config/
│   └── gophertube.yaml.example
├── man/
│   └── gophertube.1
└── internal/
    ├── app/
    ├── components/
    ├── services/
    ├── types/
    ├── interfaces/
    ├── constants/
    ├── utils/
    └── errors/
```

---

## License

[![GPL v3](https://www.gnu.org/graphics/gplv3-127x51.png)](LICENSE)

[GPL v3](https://www.gnu.org/licenses/gpl-3.0). See [LICENSE](LICENSE).

---

## Contributing

PRs and issues welcome. See [GitHub Issues](https://github.com/KrishnaSSH/GopherTube/issues) and [Pull Requests](https://github.com/KrishnaSSH/GopherTube/pulls). 