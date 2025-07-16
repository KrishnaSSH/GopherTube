package main

import (
	"fmt"
	"os"

	"gophertube/internal/app"

	"github.com/spf13/cobra"
)

var version = "dev" // will be replaced at build time with -ldflags

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gophertube",
		Short: "Terminal user interface for YouTube video search and playback",
		Long: `GopherTube is a modern terminal user interface for searching and watching 
YouTube videos using yt-dlp and mpv.

FEATURES:
  • Real-time YouTube video search
  • Direct video playback with mpv
  • Keyboard navigation interface
  • Swiss design-inspired UI
  • Configurable video quality settings
  • Responsive terminal layout

REQUIREMENTS:
  • yt-dlp: YouTube video downloader
  • mpv: Media player for video playback

CONFIGURATION:
  Create ~/.config/gophertube/gophertube.yaml for custom settings.

EXAMPLES:
  gophertube                    # Start the application
  gophertube --help            # Show this help message

For more information, see the man page: man gophertube`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			app := app.NewApp()
			if err := app.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.SetVersionTemplate("GopherTube version: {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
