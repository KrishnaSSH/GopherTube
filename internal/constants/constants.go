package constants

// UI Constants
const (
	// Video playback
	DefaultVideoQuality     = "best[height<=1080]/best"
	FallbackQuality         = "best"
	PlaybackTimeout         = 10 // Timeout for video URL extraction
	FallbackPlaybackTimeout = 8
)

// Messages
const (
	SearchingMessage   = "Searching..."
	NoVideosMessage    = "No videos found. Try searching for something!"
	LoadingMoreMessage = "Loading more videos..."
)

// Help text
const (
	VideoHelpText = "↑/↓: Move  |  Enter: Play  |  Esc: Back"
)
