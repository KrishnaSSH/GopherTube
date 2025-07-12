package services

import (
	"fmt"
	"os/exec"

	"gophertube/internal/errors"
	"gophertube/internal/utils"
)

type MPVService struct {
	config *Config
}

func NewMPVService(config *Config) *MPVService {
	return &MPVService{
		config: config,
	}
}

func (m *MPVService) PlayVideo(videoURL string) error {
	// Validate URL
	if !utils.ValidateURL(videoURL) {
		return errors.NewPlaybackError(videoURL, fmt.Errorf("invalid URL"))
	}

	// Check if mpv is available
	if err := utils.CheckCommandExists(m.config.MPVPath); err != nil {
		return errors.NewPlaybackError(videoURL, err)
	}

	// Use mpv with configurable quality settings
	formatArg := "--ytdl-format=best[height<=1080]/best[height<=720]/best" // Default: prioritize 1080p
	if m.config.VideoQuality != "best" {
		formatArg = fmt.Sprintf("--ytdl-format=%s", m.config.VideoQuality)
	}
	
	args := []string{
		formatArg,
		"--no-cache",                      // Disable cache for faster startup
		"--no-config",                     // Skip config loading
		videoURL,
	}

	cmd := exec.Command(m.config.MPVPath, args...)
	if err := cmd.Start(); err != nil {
		// Fallback to simple mpv if optimized version fails
		cmd := exec.Command(m.config.MPVPath, videoURL)
		if err := cmd.Start(); err != nil {
			return errors.NewPlaybackError(videoURL, err)
		}
	}
	return nil
}

func (m *MPVService) DownloadVideo(videoURL, outputPath string) error {
	// Validate URL
	if !utils.ValidateURL(videoURL) {
		return errors.NewPlaybackError(videoURL, fmt.Errorf("invalid URL"))
	}

	// Check if yt-dlp is available
	if err := utils.CheckCommandExists(m.config.YTDlpPath); err != nil {
		return errors.NewYTDlpError("yt-dlp", err)
	}

	args := []string{
		"--no-playlist",
		"--format", "best",
		"--output", outputPath,
		videoURL,
	}

	cmd := exec.Command(m.config.YTDlpPath, args...)
	if err := cmd.Run(); err != nil {
		return errors.NewYTDlpError("download", err)
	}

	return nil
}

func (m *MPVService) GetVideoInfo(videoURL string) (string, error) {
	// Validate URL
	if !utils.ValidateURL(videoURL) {
		return "", errors.NewPlaybackError(videoURL, fmt.Errorf("invalid URL"))
	}

	// Check if yt-dlp is available
	if err := utils.CheckCommandExists(m.config.YTDlpPath); err != nil {
		return "", errors.NewYTDlpError("yt-dlp", err)
	}

	args := []string{
		"--no-playlist",
		"--dump-json",
		videoURL,
	}

	cmd := exec.Command(m.config.YTDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewYTDlpError("info", err)
	}

	return string(output), nil
}
