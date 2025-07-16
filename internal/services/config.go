package services

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultMPVPath      = "mpv"
	DefaultYTDlpPath    = "yt-dlp"
	DefaultVideoQuality = "best[height<=1080]/best"
	DefaultSearchLimit  = 20
)

type Config struct {
	MPVPath      string `mapstructure:"mpv_path"`
	YTDlpPath    string `mapstructure:"ytdlp_path"`
	VideoQuality string `mapstructure:"video_quality"`
	DownloadPath string `mapstructure:"download_path"`
	SearchLimit  int    `mapstructure:"search_limit"`
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("gophertube")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.config/gophertube")
	v.AddConfigPath(".")

	// Set defaults
	v.SetDefault("mpv_path", DefaultMPVPath)
	v.SetDefault("ytdlp_path", DefaultYTDlpPath)
	v.SetDefault("video_quality", DefaultVideoQuality)
	v.SetDefault("download_path", getDefaultDownloadPath())
	v.SetDefault("search_limit", DefaultSearchLimit)

	_ = v.ReadInConfig() // Ignore error: config file is optional

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Defensive: ensure DownloadPath is set
	if cfg.DownloadPath == "" {
		cfg.DownloadPath = getDefaultDownloadPath()
	}
	if cfg.SearchLimit <= 0 {
		cfg.SearchLimit = DefaultSearchLimit
	}

	return &cfg, nil
}

func getDefaultDownloadPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./downloads"
	}
	return filepath.Join(homeDir, "Videos", "gophertube")
}
