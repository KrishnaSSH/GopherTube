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
	SearchLimit int `toml:"search_limit"`
}

func loadConfig() int {
	cfg := Config{SearchLimit: 30}
	data, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.config/gophertube/gophertube.toml"))
	if err == nil {
		toml.Unmarshal(data, &cfg)
	}
	if cfg.SearchLimit <= 0 {
		cfg.SearchLimit = 30
	}
	return cfg.SearchLimit
}

func printBanner() {
	fmt.Print("\033[2J\033[H")
	fmt.Println()
	fmt.Println("    \033[1;33mGopherTube\033[0m")
	fmt.Println("    \033[0;37mversion " + version + "\033[0m")
	fmt.Println()
	fmt.Println("    \033[1;36mFast Youtube Terminal UI \033[0m")
	fmt.Println("    \033[0;37mPress Ctrl+C or Esc to exit\033[0m")
	fmt.Println()
}

func printSearchPrompt(query string) {
	fmt.Print("\033[2K\r")
	fmt.Print("    \033[1;32mSearch\033[0m ")
	fmt.Print("\033[1;37m" + query + "\033[0m")
	fmt.Print("\033[1;30m█\033[0m")
}

func printProgressBar(current, total int) {
	width := 35
	filled := (current * width) / total
	bar := "\033[1;32m"
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "\033[0m"

	var phase string
	if current == 0 {
		phase = "Searching YouTube..."
	} else if current == 1 {
		phase = "Parsing search results..."
	} else if current == 2 {
		phase = "Downloading thumbnails..."
	} else {
		phase = fmt.Sprintf("Downloading thumbnails (%d/%d)", current-2, total-2)
	}

	fmt.Printf("\033[2K\r    %s %s", bar, phase)
}

func printSearchStats(videos []types.Video) {
	if len(videos) == 0 {
		return
	}

	channels := make(map[string]int)

	for _, v := range videos {
		channels[v.Author]++
	}

	fmt.Println("    \033[1;36mSearch Statistics:\033[0m")
	fmt.Printf("    \033[0;37m• Total videos: %d\033[0m\n", len(videos))
	fmt.Printf("    \033[0;37m• Unique channels: %d\033[0m\n", len(channels))

	if len(videos) > 0 {
		fmt.Printf("    \033[0;37m• Sample duration: %s\033[0m\n", videos[0].Duration)
	}

	fmt.Println()
}

func readQuery() (string, bool) {
	printBanner()
	fmt.Print("    \033[1;32mSearch\033[0m ")

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

	for {
		query, esc := readQuery()
		if esc || query == "" {
			fmt.Print("\033[2J\033[H")
			return
		}

		limit := loadConfig()
		videos, err := services.SearchYouTube(query, limit, printProgressBar)
		fmt.Println()
		fmt.Println()

		if err != nil || len(videos) == 0 {
			fmt.Println("    \033[1;31mNo results found.\033[0m")
			fmt.Println()
			fmt.Println("    \033[0;37mPress any key to search again...\033[0m")
			os.Stdin.Read(make([]byte, 1))
			continue
		}

		fmt.Printf("    \033[1;32mFound %d results!\033[0m\n", len(videos))
		printSearchStats(videos)
		time.Sleep(500 * time.Millisecond)

		selected := runFzf(videos, limit, query)
		if selected == -2 {
			continue // go back to search
		}
		if selected < 0 || selected >= len(videos) {
			continue
		}

		fmt.Printf("    \033[1;33mPlaying: %s\033[0m\n", videos[selected].Title)
		fmt.Println()
		mpvPath := "mpv"
		exec.Command(mpvPath, videos[selected].URL).Run()
	}
}

func runFzf(videos []types.Video, limit int, query string) int {
	filter := ""
	for {
		var input bytes.Buffer
		for i, v := range videos {
			thumbPath := v.ThumbnailPath
			thumbPath = strings.ReplaceAll(thumbPath, "'", "'\\''")
			fmt.Fprintf(&input, "%d\t%s\t%s\t%s\t%s\t%s\n", i, v.Title, thumbPath, v.Duration, v.Author, v.Views)
		}
		fzfArgs := []string{
			"--with-nth=2..2",
			"--delimiter=\t",
			fmt.Sprintf("--header=\033[1;36m↑/↓\033[0m to move • \033[1;33mtype\033[0m to search • \033[1;32mEnter\033[0m to select • \033[1;35mTab\033[0m to load more • \033[1;37m%d results\033[0m", len(videos)),
			"--expect=tab",
			"--bind=esc:abort",
			"--preview",
			`thumbfile={3}; w=$((FZF_PREVIEW_COLUMNS * 9 / 10)); h=$((FZF_PREVIEW_LINES * 3 / 5)); if [ -s "$thumbfile" ] && [ -f "$thumbfile" ]; then chafa --size=${w}x${h} "$thumbfile" 2>/dev/null || echo "No image preview available"; else echo "No thumbnail available"; fi; echo; title=$(printf %s {2} | fold -s -w $w | sed "s/^'//;s/'$//"); title_lines=$(echo "$title" | head -n2); if [ "$(echo "$title" | wc -l)" -gt 2 ]; then title_lines="$title_lines\n..."; fi; echo -e "\033[1;36m$title_lines\033[0m"; echo -e "\033[33mDuration:\033[0m $(echo {4} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[32mAuthor:\033[0m $(echo {5} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[35mViews:\033[0m $(echo {6} | sed s/^\'// | sed s/\'$//)"`,
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
			moreVideos, err := services.SearchYouTube(query, limit, printProgressBar)
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
