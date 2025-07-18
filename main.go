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
	"strings"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var version = "dev"

type Config struct {
	SearchLimit int `yaml:"search_limit"`
}

func loadConfig() int {
	cfg := Config{SearchLimit: 30}
	data, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.config/gophertube/gophertube.yaml"))
	if err == nil {
		yaml.Unmarshal(data, &cfg)
	}
	if cfg.SearchLimit <= 0 {
		cfg.SearchLimit = 30
	}
	return cfg.SearchLimit
}

func readQuery() (string, bool) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", false
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var query []rune
	buf := make([]byte, 1)
	for {
		fmt.Print("\033[2J\033[H> " + string(query))
		os.Stdout.Sync()
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			break
		}
		if buf[0] == 27 { // ESC
			fmt.Print("\033[2J\033[H")
			return "", true
		}
		if buf[0] == 13 { // Enter
			break
		}
		if buf[0] == 127 && len(query) > 0 { // Backspace
			query = query[:len(query)-1]
		} else if buf[0] >= 32 && buf[0] < 127 {
			query = append(query, rune(buf[0]))
		}
	}
	return string(query), false
}

func printHelp() {
	fmt.Print(`GopherTube - Terminal YouTube Search & Play

Usage:
  gophertube [options]

Options:
  -h, --help      Show this help message and exit
  -v, --version   Show version and exit
`)
}

func printVersion() {
	fmt.Println("GopherTube version", version)
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
	for {
		query, esc := readQuery()
		if esc || query == "" {
			fmt.Print("\033[2J\033[H")
			return
		}
		limit := loadConfig()
		videos, err := services.SearchYouTube(query, limit)
		if err != nil || len(videos) == 0 {
			fmt.Println("No results.")
			continue
		}
		selected := runFzf(videos, limit, query)
		if selected == -2 {
			continue // go back to search
		}
		if selected < 0 || selected >= len(videos) {
			continue
		}
		mpvPath := "mpv"
		exec.Command(mpvPath, videos[selected].URL).Run()
	}
}

func runFzf(videos []types.Video, limit int, query string) int {
	filter := ""
	for {
		var input bytes.Buffer
		for i, v := range videos {
			thumbPath := strings.ReplaceAll(v.ThumbnailPath, "'", "'\\''")
			fmt.Fprintf(&input, "%d\t%s\t%s\t%s\t%s\t%s\n", i, v.Title, thumbPath, v.Duration, v.Author, v.Views)
		}
		fzfArgs := []string{
			"--with-nth=2..2",
			"--delimiter=\t",
			"--header=↑/↓ to move, type to search, Enter to select, Tab to load more",
			"--expect=tab",
			"--bind=esc:abort",
			"--preview",
			`thumbfile={3}; w=$((FZF_PREVIEW_COLUMNS * 9 / 10)); h=$((FZF_PREVIEW_LINES * 3 / 5)); if [ -s "$thumbfile" ]; then chafa --size=${w}x${h} "$thumbfile"; else echo No image preview (not cached); fi; echo; title=$(printf %s {2} | fold -s -w $w | sed "s/^'//;s/'$//"); title_lines=$(echo "$title" | head -n2); if [ "$(echo "$title" | wc -l)" -gt 2 ]; then title_lines="$title_lines\n..."; fi; echo -e "\033[1;36m$title_lines\033[0m"; echo -e "\033[33mDuration:\033[0m $(echo {4} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[32mAuthor:\033[0m $(echo {5} | sed s/^\'// | sed s/\'$//)"; echo -e "\033[35mViews:\033[0m $(echo {6} | sed s/^\'// | sed s/\'$//)"`,
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
			fmt.Println("fzf error:", err)
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
			limit += loadConfig()
			moreVideos, err := services.SearchYouTube(query, limit)
			if err != nil || len(moreVideos) == len(videos) {
				continue
			}
			videos = moreVideos
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
