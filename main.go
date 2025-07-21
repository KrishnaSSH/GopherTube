package main

import (
	"bytes"
	"flag"
	"fmt"
	"gophertube/internal/services"
	"gophertube/internal/types"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/chzyer/readline"
)

var version = "dev"

type Config struct {
	SearchLimit   int    `toml:"search_limit"`
	Quality       string `toml:"quality"`
	DownloadsPath string `toml:"downloads_path"`
}

var config Config

func loadConfig() int {
	cfg := Config{SearchLimit: 30, Quality: "720p", DownloadsPath: os.ExpandEnv("$HOME/Videos/GopherTube")}
	data, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.config/gophertube/gophertube.toml"))
	if err == nil {
		toml.Unmarshal(data, &cfg)
	}
	if cfg.SearchLimit <= 0 {
		cfg.SearchLimit = 30
	}
	if cfg.Quality == "" {
		cfg.Quality = "720p"
	}
	if cfg.DownloadsPath == "" {
		cfg.DownloadsPath = os.ExpandEnv("$HOME/Videos/GopherTube")
	}
	config = cfg
	return cfg.SearchLimit
}

func printBanner() {
	fmt.Print("\033[2J\033[H")
	fmt.Println()
	fmt.Println("    \033[1;33mGopherTube\033[0m")
	fmt.Println("    \033[0;37mversion " + version + "\033[0m")
	fmt.Println()
	fmt.Println("    \033[1;36mFast Youtube Terminal UI\033[0m")
	fmt.Println("    \033[0;37mPress Ctrl+C or Esc to exit\033[0m")
	fmt.Println()
	fmt.Println("    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println()
}

func printSearchPrompt(query string) {
	fmt.Print("\033[2K\r")
	fmt.Print("    \033[1;32m>\033[0m ")
	fmt.Print("\033[1;37m" + query + "\033[0m")
	fmt.Print("\033[1;30m█\033[0m")
}

func printProgressBar(current, total int) {
	width := 40
	filled := (current * width) / total
	percentage := (current * 100) / total

	// Create animated progress bar with original cyan color
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "\033[1;36m█" // Cyan for filled (original color)
		} else {
			bar += "\033[0;37m░"
		}
	}
	bar += "\033[0m"

	// Add spinning animation
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinners[time.Now().UnixNano()/100000000%int64(len(spinners))]

	// Format percentage with proper padding
	percentStr := fmt.Sprintf("%3d%%", percentage)

	fmt.Printf("\033[2K\r    %s %s %s", spinner, bar, percentStr)
}

func printSearchStats(videos []types.Video) {
	if len(videos) == 0 {
		return
	}

	channels := make(map[string]int)
	hasDuration := 0

	for _, v := range videos {
		channels[v.Author]++
		if v.Duration != "" {
			hasDuration++
		}
	}

	// Calculate average duration if available
	var avgDuration string
	if hasDuration > 0 {
		avgDuration = "~" + videos[0].Duration // Simple approximation
	}

	fmt.Println("    \033[1;36mSearch Statistics:\033[0m")
	fmt.Printf("    \033[0;37m• Total videos found: \033[1;32m%d\033[0m\n", len(videos))
	fmt.Printf("    \033[0;37m• Unique channels: \033[1;33m%d\033[0m\n", len(channels))

	if avgDuration != "" {
		fmt.Printf("    \033[0;37m• Average duration: \033[1;35m%s\033[0m\n", avgDuration)
	}

	// Show top channels if there are multiple
	if len(channels) > 1 && len(videos) > 3 {
		fmt.Printf("    \033[0;37m• Most active channel: \033[1;31m%s\033[0m\n", getTopChannel(channels))
	}

	fmt.Println()
}

func getTopChannel(channels map[string]int) string {
	var topChannel string
	maxCount := 0

	for channel, count := range channels {
		if count > maxCount {
			maxCount = count
			topChannel = channel
		}
	}

	if len(topChannel) > 30 {
		return topChannel[:27] + "..."
	}
	return topChannel
}

func printSearchTips() {
	tips := []string{
		"Tip: Press Tab to load more results, Esc to go back",
		"Tip: Use ↑/↓ to navigate, Enter to select, Ctrl+C to exit",
	}

	randomTip := tips[time.Now().Unix()%int64(len(tips))]
	fmt.Printf("    \033[1;33m%s\033[0m\n", randomTip)
	fmt.Println()
}

func readQuery() (string, bool) {
	printBanner()
	fmt.Print("    \033[1;32m>\033[0m ")

	// Use raw terminal mode for proper key detection
	oldState, err := readline.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", false
	}
	defer readline.Restore(int(os.Stdin.Fd()), oldState)

	var query []rune
	buf := make([]byte, 4)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			break
		}

		// Handle escape sequences
		if buf[0] == 27 {
			if n == 1 {
				// Just Escape key pressed
				return "", true
			}
			continue
		}

		// Handle special keys
		if buf[0] == 13 { // Enter
			fmt.Println()
			break
		}
		if buf[0] == 3 { // Ctrl+C
			return "", true
		}
		if buf[0] == 127 && len(query) > 0 { // Backspace
			query = query[:len(query)-1]
			fmt.Print("\033[D \033[D")
		} else if buf[0] >= 32 && buf[0] < 127 { // Printable characters
			query = append(query, rune(buf[0]))
			fmt.Printf("%c", buf[0])
		}
	}

	return string(query), false
}

func printHelp() {
	fmt.Print("\033[1;36mGopherTube - Terminal YouTube Search & Play\033[0m\n\n")
	fmt.Print("\033[1;33mUsage:\033[0m\n")
	fmt.Print("  gophertube [options]\n\n")
	fmt.Print("\033[1;33mOptions:\033[0m\n")
	fmt.Print("  -h, --help      Show this help message and exit\n")
	fmt.Print("  -v, --version   Show version and exit\n\n")
	fmt.Print("\033[1;33mControls:\033[0m\n")
	fmt.Print("  • Type your search query and press Enter\n")
	fmt.Print("  • Use ↑/↓ to navigate results\n")
	fmt.Print("  • Press Tab to load more results\n")
	fmt.Print("  • Press Esc to go back or exit\n")
}

func printVersion() {
	fmt.Printf("\033[1;36mGopherTube\033[0m version \033[1;33m%s\033[0m\n", version)
}

func main() {
	// Parse flags
	help := false
	ver := false
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.BoolVar(&ver, "version", false, "show version")
	flag.Parse()
	if help {
		printHelp()
		return
	}
	if ver {
		printVersion()
		return
	}

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\033[1;33mExiting...\033[0m")
		os.Exit(0)
	}()

	loadConfig()
	for {
		mainMenu := []string{"Search YouTube", "Search Downloads"}
		cmd := exec.Command("fzf", "--prompt=Select mode: ")
		cmd.Stdin = strings.NewReader(strings.Join(mainMenu, "\n"))
		out, _ := cmd.Output()
		choice := strings.TrimSpace(string(out))
		if choice == "Search YouTube" {
			gophertubeYouTubeMode()
		} else if choice == "Search Downloads" {
			gophertubeDownloadsMode()
		} else {
			return
		}
	}
}

func gophertubeYouTubeMode() {
	query, esc := readQuery()
	if esc || query == "" {
		fmt.Print("\033[2J\033[H")
		return
	}

	limit := loadConfig()

	// Spinner/progress state
	progressCurrent := 0
	progressTotal := 1
	progressDone := make(chan struct{})

	// Start spinner goroutine
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

	videos, err := services.SearchYouTube(query, limit, func(current, total int) {
		progressCurrent = current
		progressTotal = total
	})

	close(progressDone)
	fmt.Print("\033[2K\r\n") // Clear progress bar/spinner line
	fmt.Println()
	fmt.Println()

	if err != nil || len(videos) == 0 {
		fmt.Println("    \033[1;31mNo results found.\033[0m")
		fmt.Println()
		fmt.Println("    \033[0;37mPress any key to search again...\033[0m")
		os.Stdin.Read(make([]byte, 1))
		return
	}

	fmt.Printf("    \033[1;32mFound %d results!\033[0m\n", len(videos))
	printSearchStats(videos)
	printSearchTips()
	// Reduced delay for faster response
	time.Sleep(200 * time.Millisecond)

	selected := runFzf(videos, limit, query)
	if selected == -2 {
		gophertubeYouTubeMode() // go back to search
	}
	if selected < 0 || selected >= len(videos) {
		gophertubeYouTubeMode()
	}

	// Show Watch/Download menu
	menu := []string{"Watch", "Download"}
	cmd := exec.Command("fzf", "--prompt=Action: ")
	cmd.Stdin = strings.NewReader(strings.Join(menu, "\n"))
	out, _ := cmd.Output()
	choice := strings.TrimSpace(string(out))
	if choice == "Download" {
		qualities := []string{"1080p", "720p", "480p", "360p"}
		cmdQ := exec.Command("fzf", "--prompt=Quality: ")
		cmdQ.Stdin = strings.NewReader(strings.Join(qualities, "\n"))
		outQ, _ := cmdQ.Output()
		selectedQ := strings.TrimSpace(string(outQ))
		if selectedQ != "" {
			// Real download logic
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
			}
			os.MkdirAll(config.DownloadsPath, 0755)
			// Sanitize filename
			filename := strings.ReplaceAll(videos[selected].Title, " ", "_")
			filename = strings.Map(func(r rune) rune {
				if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", r) {
					return r
				}
				return '_'
			}, filename)
			outputPath := fmt.Sprintf("%s/%s.%%(ext)s", config.DownloadsPath, filename)
			fmt.Printf("    \033[1;32mDownloading '%s' as %s...\033[0m\n", videos[selected].Title, selectedQ)
			ytDlpArgs := []string{"-f", format, "-o", outputPath, "--write-info-json", "--write-thumbnail", "--convert-thumbnails", "jpg", videos[selected].URL}
			cmdDl := exec.Command("yt-dlp", ytDlpArgs...)
			cmdDl.Stdout = os.Stdout
			cmdDl.Stderr = os.Stderr
			err := cmdDl.Run()
			if err == nil {
				fmt.Printf("    \033[1;32mDownload complete!\033[0m\n")
				fmt.Printf("    \033[0;37mSaved to: %s\033[0m\n", config.DownloadsPath)
			} else {
				fmt.Printf("    \033[1;31mDownload failed!\033[0m\n")
			}
			fmt.Println("    \033[0;37mPress any key to return...\033[0m")
			os.Stdin.Read(make([]byte, 1))
		}
		gophertubeYouTubeMode()
		return
	}

	// Watch as before
	fmt.Printf("    \033[1;33mPlaying: %s\033[0m\n", videos[selected].Title)
	fmt.Printf("    \033[0;37mChannel: %s\033[0m\n", videos[selected].Author)
	fmt.Printf("    \033[0;37mDuration: %s\033[0m\n", videos[selected].Duration)
	fmt.Printf("    \033[0;36mPublished: %s\033[0m\n", videos[selected].Published)
	fmt.Println()
	fmt.Println("    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println()
	mpvPath := "mpv"
	quality := config.Quality
	var mpvArgs []string
	if quality != "" {
		// Map quality to ytdl-format string
		var format string
		switch quality {
		case "1080p":
			format = "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
		case "720p":
			format = "bestvideo[height<=720]+bestaudio/best[height<=720]"
		case "480p":
			format = "bestvideo[height<=480]+bestaudio/best[height<=480]"
		case "360p":
			format = "bestvideo[height<=360]+bestaudio/best[height<=360]"
		default:
			format = "best"
		}
		mpvArgs = append(mpvArgs, "--ytdl-format="+format)
	}
	mpvArgs = append(mpvArgs, videos[selected].URL)
	exec.Command(mpvPath, mpvArgs...).Run()
}

func gophertubeDownloadsMode() {
	files, err := ioutil.ReadDir(config.DownloadsPath)
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
	fzfPreview := fmt.Sprintf(`env file={} base="%s/${file%%.*}" thumb="$base.jpg" w=$((FZF_PREVIEW_COLUMNS * 9 / 10)) h=$((FZF_PREVIEW_LINES * 3 / 5)) sh -c '[ -f "$thumb" ] && chafa --size=${w}x${h} "$thumb" 2>/dev/null || echo "No image preview available" || echo "No thumbnail available"'`, config.DownloadsPath)
	cmd := exec.Command("fzf", "--prompt=Downloads: ", "--preview", fzfPreview)
	cmd.Stdin = strings.NewReader(strings.Join(videoFiles, "\n"))
	out, _ := cmd.Output()
	selected := strings.TrimSpace(string(out))
	if selected == "" {
		return
	}
	filePath := config.DownloadsPath + "/" + selected
	fmt.Printf("    \033[1;33mPlaying: %s\033[0m\n", selected)
	fmt.Println()
	fmt.Println("    \033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println()
	mpvPath := "mpv"
	exec.Command(mpvPath, filePath).Run()
}

func runFzf(videos []types.Video, limit int, query string) int {
	filter := ""
	for {
		var input bytes.Buffer
		for i, v := range videos {
			thumbPath := v.ThumbnailPath
			thumbPath = strings.ReplaceAll(thumbPath, "'", "'\\''")
			fmt.Fprintf(&input, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", i, v.Title, thumbPath, v.Duration, v.Author, v.Views, v.Description, v.Published)
		}
		fzfArgs := []string{
			"--with-nth=2..2",
			"--delimiter=\t",
			fmt.Sprintf("--header=\033[1;36m↑/↓\033[0m to move • \033[1;33mtype\033[0m to search • \033[1;32mEnter\033[0m to select • \033[1;35mTab\033[0m to load more • \033[1;37m%d results • \033[1;35m%s\033[0m", len(videos), query),
			"--expect=tab",
			"--bind=esc:abort",
			"--border=rounded",
			"--margin=1,1",
			"--preview",
			`env thumbfile={3} w=$((FZF_PREVIEW_COLUMNS * 9 / 10)) h=$((FZF_PREVIEW_LINES * 3 / 5)) sh -c '[ -s "$thumbfile" ] && [ -f "$thumbfile" ] && chafa --size=${w}x${h} "$thumbfile" 2>/dev/null || echo "No image preview available" || echo "No thumbnail available"'; echo; echo -e "\033[33mDuration:\033[0m $(echo {4} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[36mPublished:\033[0m $(echo {8} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[32mAuthor:\033[0m $(echo {5} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[35mViews:\033[0m $(echo {6} | sed s/^\'// | sed s/\'$//)"`,
		}
		if filter != "" {
			fzfArgs = append(fzfArgs, "--query="+filter)
		}
		cmd := exec.Command("fzf", fzfArgs...)
		cmd.Stdin = &input
		pr, pw, _ := os.Pipe()
		cmd.Stdout = pw
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Println("\033[1;31mfzf error:\033[0m", err)
			return -1
		}
		pw.Close()
		selected, _ := io.ReadAll(pr)
		cmd.Wait()
		lines := strings.Split(strings.TrimSpace(string(selected)), "\n")
		if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
			return -2 // user pressed escape in fzf
		}
		if lines[0] == "tab" {
			fmt.Printf("    \033[1;35mLoading more results...\033[0m\n")
			limit += loadConfig()
			moreVideos, err := services.SearchYouTube(query, limit, func(current, total int) {
				// progressCurrent = current // This line was removed
				// progressTotal = total // This line was removed
			})
			if err != nil || len(moreVideos) == len(videos) {
				continue
			}
			videos = moreVideos
			fmt.Printf("    \033[1;32mLoaded %d total results!\033[0m\n", len(videos))
			printSearchStats(videos)
			continue
		}
		line := lines[0]
		if line == "" {
			return -1
		}
		idxStr := strings.SplitN(line, "\t", 2)[0]
		idx := 0
		fmt.Sscanf(idxStr, "%d", &idx)
		return idx
	}
}
