# GopherTube

A terminal user interface for searching and watching YouTube videos using yt-dlp and mpv.

## Features

- Terminal-based user interface with Swiss design principles
- YouTube video search with real-time results
- Direct video playback with mpv
- Keyboard navigation for video browsing
- Responsive layout that adapts to terminal size
- YAML-based configuration system

## Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **yt-dlp** - YouTube downloader
  ```bash
  pip install yt-dlp
  ```
- **mpv** - Media player
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

- **Bubble Tea** for TUI framework
- **Lip Gloss** for styling
- **Cobra** for CLI commands
- **Viper** for configuration management

## Troubleshooting

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

This project is licensed under the [GNU GPLv3](https://github.com/KrishnaSSH/GopherTube/blob/main/LICENSE). 