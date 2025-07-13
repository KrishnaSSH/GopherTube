package services

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"gophertube/internal/constants"
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

	// Check if yt-dlp and mpv are available
	if err := utils.CheckCommandExists(m.config.YTDlpPath); err != nil {
		return errors.NewPlaybackError(videoURL, fmt.Errorf("yt-dlp not found: %w", err))
	}
	if err := utils.CheckCommandExists(m.config.MPVPath); err != nil {
		return errors.NewPlaybackError(videoURL, err)
	}

	// Use yt-dlp to get direct video URL
	cmd := exec.Command(m.config.YTDlpPath,
		"-f", constants.DefaultVideoQuality,
		"--get-url",
		"--no-playlist",
		"--no-warnings",
		"--quiet",
		"--no-check-certificates",
		"--no-cache-dir",
		videoURL,
	)

	// Add timeout for faster failure
	ctx, cancel := context.WithTimeout(context.Background(), constants.PlaybackTimeout*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		// Fallback to simplest format for maximum compatibility
		cmd = exec.Command(m.config.YTDlpPath,
			"-f", constants.FallbackQuality,
			"--get-url",
			"--no-playlist",
			"--no-warnings",
			"--quiet",
			"--no-check-certificates",
			"--no-cache-dir",
			videoURL,
		)

		ctx2, cancel2 := context.WithTimeout(context.Background(), constants.FallbackPlaybackTimeout*time.Second)
		defer cancel2()
		cmd = exec.CommandContext(ctx2, cmd.Path, cmd.Args[1:]...)

		output, err = cmd.Output()
		if err != nil || len(output) == 0 {
			return errors.NewPlaybackError(videoURL, fmt.Errorf("yt-dlp failed to get video URL: %w", err))
		}
	}

	// yt-dlp may return one or two URLs (video and audio)
	urls := []string{}
	for _, line := range splitLines(string(output)) {
		if line != "" {
			urls = append(urls, line)
		}
	}
	if len(urls) == 0 {
		return errors.NewPlaybackError(videoURL, fmt.Errorf("no playable URL found"))
	}

	// Pass the URLs to mpv
	mpvArgs := append([]string{
		"--no-config",
		"--no-cache",
		"--no-ytdl",
		"--geometry=50%:50%",
	}, urls...)

	cmd = exec.Command(m.config.MPVPath, mpvArgs...)
	if err := cmd.Start(); err != nil {
		return errors.NewPlaybackError(videoURL, err)
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
		"--no-warnings",
		"--quiet",
		"--no-check-certificates",
		"--no-cache-dir",
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
		"--no-warnings",
		"--quiet",
		"--no-check-certificates",
		"--no-cache-dir",
		videoURL,
	}

	cmd := exec.Command(m.config.YTDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewYTDlpError("info", err)
	}

	return string(output), nil
}

// splitLines splits a string into lines, handling both \n and \r\n
func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
