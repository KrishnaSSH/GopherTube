package constants

// UI Constants
const (
	AppTitle = "GopherTube"

	// Colors
	PrimaryColor   = "#00ff00"
	SecondaryColor = "#888888"
	ErrorColor     = "#ff0000"

	// Responsive Design
	MinTerminalWidth  = 40
	MinTerminalHeight = 10
	MaxTerminalWidth  = 200
	MaxTerminalHeight = 100

	// Spacing
	DefaultPadding = 2
	DefaultMargin  = 4

	// Search
	MaxSearchResults = 8 // Reduced from 10 for speed
	SearchTimeout    = 8 // Reduced from 30 seconds
	FallbackTimeout  = 5 // Fallback timeout for search

	// Video playback
	DefaultVideoQuality     = "best[height<=1080]/best"
	FallbackQuality         = "best"
	PlaybackTimeout         = 10 // Timeout for video URL extraction
	FallbackPlaybackTimeout = 8

	// Commands
	MPVCommand   = "mpv"
	YTDlpCommand = "yt-dlp"
)

// Messages
const (
	SearchingMessage    = "Searching..."
	NoVideosMessage     = "No videos found. Try searching for something!"
	ReadyMessage        = "Ready"
	ErrorPrefix         = "Error: "
	SearchFailedMessage = "Search failed: "
	FoundVideosMessage  = "Found %d videos"
	LoadingMoreMessage  = "Loading more videos..."
)

// Help text
const (
	SearchHelpText = "Enter: Search  |  Ctrl+C or Esc: Quit"
	VideoHelpText  = "↑/↓: Move  |  Enter: Play  |  Esc: Back"

	// Compact help text for small terminals
	CompactSearchHelpText = "Enter: Search  |  Esc: Quit"
	CompactVideoHelpText  = "↑/↓: Move  |  Enter: Play  |  Esc: Back"
)

// File paths
const (
	ConfigFileName = "gophertube.yaml"
	ConfigDir      = "$HOME/.config/gophertube"
	DownloadsDir   = "~/Videos/gophertube"
)

// Performance settings
const (
	// Cache settings
	MaxCacheSize = 100 // Maximum number of cached search results

	// yt-dlp optimization flags
	YTDlpOptimizedFlags = "--no-warnings --quiet --no-check-certificates --no-cache-dir"

	// MPV optimization flags
	MPVOptimizedFlags = "--no-config --no-cache --no-ytdl --no-video-title-show"
)
