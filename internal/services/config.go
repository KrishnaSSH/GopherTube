package services

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	MPVPath      string
	YTDlpPath    string
	VideoQuality string
	DownloadPath string
	SearchLimit  int
}

func NewConfig() *Config {
	config := &Config{
		MPVPath:      "mpv",
		YTDlpPath:    "yt-dlp",
		VideoQuality: "best[height<=1080]/best",
		DownloadPath: getDefaultDownloadPath(),
		SearchLimit:  20, // changed from 8 to 20
	}

	// Load config from file
	viper.SetConfigName("gophertube")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/gophertube")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		if mpvPath := viper.GetString("mpv_path"); mpvPath != "" {
			config.MPVPath = mpvPath
		}
		if ytdlpPath := viper.GetString("ytdlp_path"); ytdlpPath != "" {
			config.YTDlpPath = ytdlpPath
		}
		if quality := viper.GetString("video_quality"); quality != "" {
			config.VideoQuality = quality
		}
		if downloadPath := viper.GetString("download_path"); downloadPath != "" {
			config.DownloadPath = downloadPath
		}
		if searchLimit := viper.GetInt("search_limit"); searchLimit > 0 {
			config.SearchLimit = searchLimit
		}
	}

	return config
}

// ConfigService interface implementation
func (c *Config) GetMPVPath() string {
	return c.MPVPath
}

func (c *Config) GetYTDlpPath() string {
	return c.YTDlpPath
}

func (c *Config) GetVideoQuality() string {
	return c.VideoQuality
}

func (c *Config) GetDownloadPath() string {
	return c.DownloadPath
}

func (c *Config) GetSearchLimit() int {
	return c.SearchLimit
}

func getDefaultDownloadPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./downloads"
	}
	return filepath.Join(homeDir, "Videos", "gophertube")
}
