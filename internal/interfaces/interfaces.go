package interfaces

import (
	"gophertube/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

// VideoService defines the interface for video operations
type VideoService interface {
	SearchVideos(query string, maxResults int) ([]types.Video, error)
	GetVideoInfo(videoURL string) (string, error)
}

// PlayerService defines the interface for video playback
type PlayerService interface {
	PlayVideo(videoURL string) error
	DownloadVideo(videoURL, outputPath string) error
}

// ConfigService defines the interface for configuration management
type ConfigService interface {
	GetMPVPath() string
	GetYTDlpPath() string
	GetVideoQuality() string
	GetDownloadPath() string
}

// UIComponent defines the interface for UI components
type UIComponent interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	SetSize(width, height int)
}
