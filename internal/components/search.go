package components

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gophertube/internal/errors"
	"gophertube/internal/utils"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchComponent struct {
	textInput textinput.Model
	spinner   spinner.Model
	width     int
	height    int
	isLoading bool
	query     string
}

type SearchResultMsg struct {
	Videos []Video
	Error  string
}

func NewSearchComponent() *SearchComponent {
	ti := textinput.New()
	ti.Placeholder = "Search for YouTube videos..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &SearchComponent{
		textInput: ti,
		spinner:   s,
		width:     80,
		height:    20,
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
			Render(s.spinner.View() + " Searching...")
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

func (s *SearchComponent) SearchWithQuery(query string) ([]Video, error) {
	return s.searchWithYtDlp(query)
}

func (s *SearchComponent) searchVideos(query string) tea.Cmd {
	return func() tea.Msg {
		// Use yt-dlp to search for videos
		videos, err := s.searchWithYtDlp(query)
		if err != nil {
			return SearchResultMsg{Error: err.Error()}
		}

		return SearchResultMsg{Videos: videos}
	}
}

func (s *SearchComponent) searchWithYtDlp(query string) ([]Video, error) {
	// Check if yt-dlp is available
	if err := utils.CheckCommandExists("yt-dlp"); err != nil {
		return nil, errors.NewYTDlpError("yt-dlp", err)
	}

	// Use yt-dlp to search YouTube with faster options
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--flat-playlist",
		"--max-downloads", "10", // Get 10 videos
		"--no-warnings",                     // Reduce output
		"--quiet",                           // Quiet mode for speed
		fmt.Sprintf("ytsearch10:%s", query), // Get 10 search results
	)

	// Add timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	output, err := cmd.Output()
	if err != nil {
		// Try with fewer results if the first one fails
		cmd = exec.Command("yt-dlp",
			"--dump-json",
			"--no-playlist",
			"--max-downloads", "5",
			"--no-warnings",
			"--quiet",
			fmt.Sprintf("ytsearch5:%s", query),
		)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel2()
		cmd = exec.CommandContext(ctx2, cmd.Path, cmd.Args[1:]...)

		output, err = cmd.Output()
		if err != nil {
			return nil, errors.NewYTDlpError("search", err)
		}
	}

	// Parse the JSON output
	var videos []Video
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var videoData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &videoData); err != nil {
			continue // Skip invalid JSON lines
		}

		video := Video{
			Title:       utils.SafeGetString(videoData, "title"),
			Author:      utils.SafeGetString(videoData, "uploader"),
			Duration:    utils.FormatDuration(utils.SafeGetInt(videoData, "duration")),
			Views:       utils.FormatViews(utils.SafeGetInt(videoData, "view_count")),
			URL:         utils.SafeGetString(videoData, "webpage_url"),
			Thumbnail:   utils.SafeGetString(videoData, "thumbnail"),
			Description: utils.SafeGetString(videoData, "description"),
		}

		// Only add videos with valid URLs
		if video.URL != "" {
			videos = append(videos, video)
		}
	}

	if len(videos) == 0 {
		return nil, errors.NewSearchError(query, fmt.Errorf("no videos found"))
	}

	return videos, nil
}
 