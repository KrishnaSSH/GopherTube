package constants

// UI Constants
const (
	AppTitle = "GopherTube"

	// Colors
	PrimaryColor   = "#00ff00"
	SecondaryColor = "#888888"
	ErrorColor     = "#ff0000"

	// Spacing
	DefaultPadding = 2
	DefaultMargin  = 4

	// Search
	MaxSearchResults = 8
	SearchTimeout    = 8
	FallbackTimeout  = 5

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
)

// File paths
const (
	ConfigFileName = "gophertube.yaml"
	ConfigDir      = "$HOME/.config/gophertube"
	DownloadsDir   = "~/Videos/gophertube"
)
