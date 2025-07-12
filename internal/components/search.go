package components

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"gophertube/internal/constants"
	"gophertube/internal/errors"
	"gophertube/internal/services"
	"gophertube/internal/types"
	"gophertube/internal/utils"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchComponent struct {
	textInput        textinput.Model
	spinner          spinner.Model
	width            int
	height           int
	isLoading        bool
	query            string
	cache            map[string][]types.Video
	cacheMux         sync.RWMutex
	invidiousService *services.InvidiousService
}

type SearchResultMsg struct {
	Videos []types.Video
	Error  string
}

func NewSearchComponent(config *services.Config) *SearchComponent {
	ti := textinput.New()
	ti.Placeholder = "Search for YouTube videos..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &SearchComponent{
		textInput:        ti,
		spinner:          s,
		width:            80,
		height:           20,
		cache:            make(map[string][]types.Video),
		invidiousService: services.NewInvidiousService(config),
	}
}

func (s *SearchComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (s *SearchComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if s.textInput.Value() != "" && !s.isLoading {
				s.query = s.textInput.Value()
				s.isLoading = true
				return s, tea.Batch(
					s.spinner.Tick,
					s.searchVideos(s.query),
				)
			}
		case "ctrl+c":
			return s, tea.Quit
		}
	case spinner.TickMsg:
		if s.isLoading {
			var cmd tea.Cmd
			s.spinner, cmd = s.spinner.Update(msg)
			return s, cmd
		}
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *SearchComponent) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Left).
		Width(s.width).
		Render("GopherTube")

	searchBox := lipgloss.NewStyle().
		Padding(0, 2).
		Width(s.width).
		Align(lipgloss.Left).
		Render(s.textInput.View())

	var content string
	if s.isLoading {
		content = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(s.width).
			Render(s.spinner.View() + " " + constants.SearchingMessage)
	} else {
		content = searchBox
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Left).
		Width(s.width).
		Render("Enter: Search  |  Ctrl+C or Esc: Quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		help,
	)
}

func (s *SearchComponent) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.textInput.Width = width - 10
}

func (s *SearchComponent) ResetLoading() {
	s.isLoading = false
}

func (s *SearchComponent) GetCurrentQuery() string {
	return s.textInput.Value()
}

func (s *SearchComponent) SearchWithQuery(query string) ([]types.Video, error) {
	return s.searchWithInvidious(query)
}

func (s *SearchComponent) searchVideos(query string) tea.Cmd {
	return func() tea.Msg {
		// Check cache first
		s.cacheMux.RLock()
		if cached, exists := s.cache[query]; exists {
			s.cacheMux.RUnlock()
			return SearchResultMsg{Videos: cached}
		}
		s.cacheMux.RUnlock()

		// Use Invidious to search for videos
		videos, err := s.searchWithInvidious(query)
		if err != nil {
			return SearchResultMsg{Error: err.Error()}
		}

		// Cache the results
		s.cacheMux.Lock()
		s.cache[query] = videos
		s.cacheMux.Unlock()

		return SearchResultMsg{Videos: videos}
	}
}

func (s *SearchComponent) searchWithInvidious(query string) ([]types.Video, error) {
	// Performance timer
	timer := utils.StartTimer("invidious_search")
	defer timer.StopTimerWithLog()

	// Use Invidious service for search
	videos, err := s.invidiousService.SearchVideos(query)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

// Fallback to yt-dlp search if Invidious fails
func (s *SearchComponent) searchWithYtDlp(query string) ([]types.Video, error) {
	// Performance timer
	timer := utils.StartTimer("ytdlp_search")
	defer timer.StopTimerWithLog()

	// Check if yt-dlp is available
	if err := utils.CheckCommandExists("yt-dlp"); err != nil {
		return nil, errors.NewYTDlpError("yt-dlp", err)
	}

	// Optimized yt-dlp search with faster parameters
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--flat-playlist",
		"--max-downloads", fmt.Sprintf("%d", constants.MaxSearchResults),
		"--no-warnings",
		"--quiet",
		"--no-check-certificates",                   // Skip SSL verification for speed
		"--no-cache-dir",                            // Disable cache for faster startup
		"--extractor-args", "youtube:skip=hls,dash", // Skip complex formats
		fmt.Sprintf("ytsearch%d:%s", constants.MaxSearchResults, query),
	)

	// Shorter timeout for faster failure
	ctx, cancel := context.WithTimeout(context.Background(), constants.SearchTimeout*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	output, err := cmd.Output()
	if err != nil {
		// Try with even fewer results and different approach
		cmd = exec.Command("yt-dlp",
			"--dump-json",
			"--no-playlist",
			"--max-downloads", "5",
			"--no-warnings",
			"--quiet",
			"--no-check-certificates",
			"--no-cache-dir",
			fmt.Sprintf("ytsearch5:%s", query),
		)

		ctx2, cancel2 := context.WithTimeout(context.Background(), constants.FallbackTimeout*time.Second)
		defer cancel2()
		cmd = exec.CommandContext(ctx2, cmd.Path, cmd.Args[1:]...)

		output, err = cmd.Output()
		if err != nil {
			return nil, errors.NewYTDlpError("search", err)
		}
	}

	// Parse the JSON output more efficiently
	var videos []types.Video
	lines := utils.FastSplit(utils.FastTrim(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var videoData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &videoData); err != nil {
			continue // Skip invalid JSON lines
		}

		// Only extract essential fields for speed
		video := types.Video{
			Title:       utils.SafeGetString(videoData, "title"),
			Author:      utils.SafeGetString(videoData, "uploader"),
			Duration:    utils.FormatDuration(utils.SafeGetInt(videoData, "duration")),
			Views:       utils.FormatViews(utils.SafeGetInt(videoData, "view_count")),
			URL:         utils.SafeGetString(videoData, "webpage_url"),
			Thumbnail:   utils.SafeGetString(videoData, "thumbnail"),
			Description: utils.SafeGetString(videoData, "description"),
		}

		// Only add videos with valid URLs and titles
		if video.URL != "" && video.Title != "" {
			videos = append(videos, video)
		}
	}

	if len(videos) == 0 {
		return nil, errors.NewSearchError(query, fmt.Errorf("no videos found"))
	}

	return videos, nil
}

// ClearCache clears the search cache
func (s *SearchComponent) ClearCache() {
	s.cacheMux.Lock()
	s.cache = make(map[string][]types.Video)
	s.cacheMux.Unlock()
}
