package app

import (
	"fmt"
	"gophertube/internal/services"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/urfave/cli/v3"
)

// MediaPlayer represents available media players
type MediaPlayer struct {
	Name string
	Path string
	Type string // "video" or "audio"
}

// checkAvailablePlayer checks for MPV first, then VLC as fallback
func checkAvailablePlayer() *MediaPlayer {
	// Prefer MPV for better performance and terminal integration
	if path, err := exec.LookPath("mpv"); err == nil {
		return &MediaPlayer{
			Name: "mpv",
			Path: path,
			Type: "video",
		}
	}
	
	// Fallback to VLC if MPV not available
	if path, err := exec.LookPath("vlc"); err == nil {
		return &MediaPlayer{
			Name: "vlc",
			Path: path,
			Type: "video",
		}
	}
	
	return nil
}

// checkAvailableAudioPlayer checks for MPV first for audio playback
func checkAvailableAudioPlayer() *MediaPlayer {
	// Prefer MPV for audio - excellent terminal integration and performance
	if path, err := exec.LookPath("mpv"); err == nil {
		return &MediaPlayer{
			Name: "mpv",
			Path: path,
			Type: "audio",
		}
	}
	
	// Fallback to VLC if MPV not available
	if path, err := exec.LookPath("vlc"); err == nil {
		return &MediaPlayer{
			Name: "vlc",
			Path: path,
			Type: "audio",
		}
	}
	
	return nil
}

// playWithPlayer plays media using the detected player
func playWithPlayer(player *MediaPlayer, url string, isAudioOnly bool) error {
	var args []string
	
	switch player.Name {
	case "mpv":
		if player.Type == "audio" || isAudioOnly {
			args = []string{
				"--no-video",
				
			}
		} else {
			args = []string{
				"--no-terminal",
				"--hwdec=auto", // Enable hardware decoding
				"--vo=gpu",     // Use GPU video output
				"--quiet",      // Reduce messages but keep controls
			}
		}
		args = append(args, url)
	
	case "vlc":
		if player.Type == "audio" || isAudioOnly {
			// VLC in audio-only mode with interface suppression
			args = []string{
				"--intf", "dummy",           // Use dummy interface (no GUI)
				"--play-and-exit",
				"--no-video",
				"--no-video-title-show",
				"--quiet",                   // Reduce verbosity
				"--no-sout-video",          // Disable video output
			}
		} else {
			// VLC in video mode with better error handling
			args = []string{
				"--play-and-exit",
				"--no-video-title-show",
				"--quiet",                   // Reduce error messages
				"--avcodec-hw=any",         // Enable hardware decoding
				"--no-audio-time-stretch",  // Prevent audio sync issues
				"--clock-jitter=0",         // Reduce timing issues
				"--network-caching=1000",   // Increase network buffer
				"--file-caching=1000",      // Increase file buffer
			}
		}
		args = append(args, url)
	}
	
	cmd := exec.Command(player.Path, args...)
	if player.Name == "vlc" {
	cmd.Stdout = nil
	cmd.Stderr = nil
} else {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

	
	return cmd.Run()
}

func gophertubeYouTubeMode(cmd *cli.Command) {
	query, esc := readQuery()
	if esc || query == "" {
		fmt.Print("\033[2J\033[H")
		return
	}
	
	progressCurrent := 0
	progressTotal := 1
	progressDone := make(chan struct{})
	go func() {
		for {
			select {
			case <-progressDone:
				return
			default:
				printProgressBar(progressCurrent, progressTotal)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	
	videos, err := services.SearchYouTube(query, cmd.Int(FlagSearchLimit), func(current, total int) {
		progressCurrent = current
		progressTotal = total
	})
	close(progressDone)
	fmt.Print("\033[2K\r\n")
	fmt.Println("\n")
	
	if err != nil || len(videos) == 0 {
		fmt.Println("    \033[1;31mNo results found.\033[0m")
		fmt.Println("\n    \033[0;37mPress any key to search again...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		return
	}
	
	fmt.Printf("    \033[1;32mFound %d results!\033[0m\n", len(videos))
	printSearchStats(videos)
	printSearchTips()
	time.Sleep(200 * time.Millisecond)
	
	selected := runFzf(videos, cmd.Int(FlagSearchLimit), query)
	if selected == -2 || selected < 0 || selected >= len(videos) {
		gophertubeYouTubeMode(cmd)
		return
	}
	
	action := exec.Command("fzf", "--prompt=Action: ")
	action.Stdin = strings.NewReader(strings.Join([]string{"Watch", "Download", "Audio"}, "\n"))
	out, _ := action.Output()
	choice := strings.TrimSpace(string(out))
	
	if choice == "Download" {
		qualities := []string{"1080p", "720p", "480p", "360p", "Audio"}
		actionQ := exec.Command("fzf", "--prompt=Quality: ")
		actionQ.Stdin = strings.NewReader(strings.Join(qualities, "\n"))
		outQ, _ := actionQ.Output()
		selectedQ := strings.TrimSpace(string(outQ))
		
		if selectedQ != "" {
			format := "best"
			switch selectedQ {
			case "1080p":
				format = "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
			case "720p":
				format = "bestvideo[height<=720]+bestaudio/best[height<=720]"
			case "480p":
				format = "bestvideo[height<=480]+bestaudio/best[height<=480]"
			case "360p":
				format = "bestvideo[height<=360]+bestaudio/best[height<=360]"
			case "Audio":
				format = "bestaudio"
			}
			
			os.MkdirAll(cmd.String(FlagDownloadsPath), 0755)
			filename := strings.ReplaceAll(videos[selected].Title, " ", "_")
			filename = strings.Map(func(r rune) rune {
				if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", r) {
					return r
				}
				return '_'
			}, filename)
			outputPath := fmt.Sprintf("%s/%s.%%(ext)s", cmd.String(FlagDownloadsPath), filename)
			
			fmt.Printf("    \033[1;32mDownloading '%s' as %s...\033[0m\n", videos[selected].Title, selectedQ)
			ytDlpArgs := []string{"-f", format, "-o", outputPath, "--write-info-json", "--write-thumbnail", "--convert-thumbnails", "jpg", videos[selected].URL}
			if format == "bestaudio" {
				ytDlpArgs = []string{"-x", "-f", format, "-o", outputPath, "--write-info-json", "--write-thumbnail", "--convert-thumbnails", "jpg", videos[selected].URL}
			}
			
			actionDl := exec.Command("yt-dlp", ytDlpArgs...)
			actionDl.Stdout = os.Stdout
			actionDl.Stderr = os.Stderr
			if err := actionDl.Run(); err == nil {
				fmt.Printf("    \033[1;32mDownload complete!\033[0m\n")
				fmt.Printf("    \033[0;37mSaved to: %s\033[0m\n", cmd.String(FlagDownloadsPath))
			} else {
				fmt.Printf("    \033[1;31mDownload failed!\033[0m\n")
			}
			fmt.Println("    \033[0;37mPress any key to return...\033[0m")
			os.Stdin.Read(make([]byte, 1))
		}
		gophertubeYouTubeMode(cmd)
		return
	}
	
	// Check for available media player
	player := checkAvailablePlayer()
	if player == nil {
		fmt.Println("    \033[1;31mNo media player found!\033[0m")
		fmt.Println("    \033[0;37mPlease install MPV (recommended) or VLC to play videos.\033[0m")
		fmt.Println("    \033[0;33mInstall MPV: sudo apt install mpv (Ubuntu) | brew install mpv (macOS)\033[0m")
		fmt.Println("    \033[0;37mPress any key to return...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		gophertubeYouTubeMode(cmd)
		return
	}
	
	// If user chooses to just play Audio
	if choice == "Audio" {
		// Check for available audio player (terminal-based preferred)
		audioPlayer := checkAvailableAudioPlayer()
		if audioPlayer == nil {
			fmt.Println("    \033[1;31mNo audio player found!\033[0m")
			fmt.Println("    \033[0;37mPlease install MPV (recommended) or VLC to play audio.\033[0m")
			fmt.Println("    \033[0;33mInstall MPV: sudo apt install mpv (Ubuntu) | brew install mpv (macOS)\033[0m")
			fmt.Println("    \033[0;37mPress any key to return...\033[0m")
			os.Stdin.Read(make([]byte, 1))
			gophertubeYouTubeMode(cmd)
			return
		}
		
		fmt.Printf("    \033[1;33mPlaying Audio with %s: %s\033[0m\n", strings.ToUpper(audioPlayer.Name), videos[selected].Title)
		fmt.Printf("    \033[0;37mChannel: %s\033[0m\n", videos[selected].Author)
		fmt.Printf("    \033[0;37mDuration: %s\033[0m\n", videos[selected].Duration)
		fmt.Printf("    \033[0;36mPublished: %s\033[0m\n", videos[selected].Published)
		fmt.Println("    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		
		// Show different controls based on player
		switch audioPlayer.Name {
		case "mpv":
			fmt.Println("    \033[0;33mControls: 'q' to quit, SPACE to pause/resume, ←→ to seek\033[0m")
		case "vlc":
			fmt.Println("    \033[0;33mControls: Ctrl+C to quit (VLC has limited terminal controls)\033[0m")
		}
		
		fmt.Println("    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
		
		// Extract direct audio stream URL with better format selection
		audioCmd := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]/bestaudio", "-g", videos[selected].URL)
		streamURLBytes, err := audioCmd.Output()
		if err != nil {
			fmt.Println("    \033[1;31mFailed to get direct audio URL.\033[0m")
			fmt.Println("    \033[0;37mMake sure yt-dlp is installed.\033[0m")
			fmt.Println("    \033[0;37mPress any key to return...\033[0m")
			os.Stdin.Read(make([]byte, 1))
			gophertubeYouTubeMode(cmd)
			return
		}
		streamURL := strings.TrimSpace(string(streamURLBytes))
		
		// Play audio with detected terminal player
		if err := playWithPlayer(audioPlayer, streamURL, true); err != nil {
			fmt.Printf("    \033[1;31mFailed to play audio with %s.\033[0m\n", audioPlayer.Name)
		}
		
		fmt.Println("    \033[0;37mPress Enter to return.\033[0m")
		os.Stdin.Read(make([]byte, 1))
		gophertubeYouTubeMode(cmd)
		return
	}
	
	fmt.Printf("    \033[1;33mPlaying with %s: %s\033[0m\n", strings.ToUpper(player.Name), videos[selected].Title)
	fmt.Printf("    \033[0;37mChannel: %s\033[0m\n", videos[selected].Author)
	fmt.Printf("    \033[0;37mDuration: %s\033[0m\n", videos[selected].Duration)
	fmt.Printf("    \033[0;36mPublished: %s\033[0m\n", videos[selected].Published)
	fmt.Println("\n    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	
	quality := cmd.String(FlagQuality)
	isAudioOnly := quality == "Audio"
	
	// Extract direct streamable URL using yt-dlp with better format selection
	var ytDlpCmd *exec.Cmd
	if isAudioOnly {
		ytDlpCmd = exec.Command("yt-dlp", "-f", "bestaudio", "-g", videos[selected].URL)
	} else {
		// Try to get the best video format that's compatible
		ytDlpCmd = exec.Command("yt-dlp", "-f", "best[height<=1080]/best", "-g", videos[selected].URL)
	}
	
	streamURLBytes, err := ytDlpCmd.Output()
	if err != nil {
		fmt.Println("    \033[1;31mFailed to get direct video URL.\033[0m")
		fmt.Println("    \033[0;37mMake sure yt-dlp is installed.\033[0m")
		return
	}
	streamURL := strings.TrimSpace(string(streamURLBytes))
	
	// Play with detected player
	if err := playWithPlayer(player, streamURL, isAudioOnly); err != nil {
		fmt.Printf("    \033[1;31mFailed to play video with %s.\033[0m\n", player.Name)
	}
}

func gophertubeDownloadsMode(cmd *cli.Command) {
	files, err := os.ReadDir(cmd.String(FlagDownloadsPath))
	if err != nil || len(files) == 0 {
		fmt.Println("    \033[1;31mNo downloaded videos found.\033[0m")
		fmt.Println("    \033[0;37mPress any key to return to main menu...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		return
	}
	
	var videoFiles []string
	for _, f := range files {
		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".mp4") || strings.HasSuffix(f.Name(), ".mkv") || strings.HasSuffix(f.Name(), ".webm") || strings.HasSuffix(f.Name(), ".avi")) {
			videoFiles = append(videoFiles, f.Name())
		}
	}
	
	if len(videoFiles) == 0 {
		fmt.Println("    \033[1;31mNo downloaded videos found.\033[0m")
		fmt.Println("    \033[0;37mPress any key to return to main menu...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		return
	}
	
	fzfPreview := fmt.Sprintf(`env file={} base="%s/${file%%.*}" thumb="$base.jpg" w=$((FZF_PREVIEW_COLUMNS * 9 / 10)) h=$((FZF_PREVIEW_LINES * 3 / 5)) sh -c '[ -f "$thumb" ] && chafa --size=${w}x${h} "$thumb" 2>/dev/null || echo "No image preview available" || echo "No thumbnail available"'`, cmd.String(FlagDownloadsPath))
	action := exec.Command("fzf", "--prompt=Downloads: ", "--preview", fzfPreview)
	action.Stdin = strings.NewReader(strings.Join(videoFiles, "\n"))
	out, _ := action.Output()
	selected := strings.TrimSpace(string(out))
	
	if selected == "" {
		return
	}
	
	// Check for available media player
	player := checkAvailablePlayer()
	if player == nil {
		fmt.Println("    \033[1;31mNo media player found!\033[0m")
		fmt.Println("    \033[0;37mPlease install MPV (recommended) or VLC to play videos.\033[0m")
		fmt.Println("    \033[0;33mInstall MPV: sudo apt install mpv (Ubuntu) | brew install mpv (macOS)\033[0m")
		fmt.Println("    \033[0;37mPress any key to return...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		return
	}
	
	filePath := cmd.String(FlagDownloadsPath) + "/" + selected
	fmt.Printf("    \033[1;33mPlaying with %s: %s\033[0m\n", strings.ToUpper(player.Name), selected)
	fmt.Println("\n    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")
	
	// Play with detected player
	if err := playWithPlayer(player, filePath, false); err != nil {
		fmt.Printf("    \033[1;31mFailed to play video with %s.\033[0m\n", player.Name)
	}
}