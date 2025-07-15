package components

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gophertube/internal/constants"
	"gophertube/internal/services"
	"gophertube/internal/types"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchComponent struct {
	textInput          textinput.Model
	spinner            spinner.Model
	width              int
	height             int
	isLoading          bool
	query              string
	cache              map[string][]types.Video
	cacheMux           sync.RWMutex
	recentSearches     []string
	lastSearchDuration time.Duration
	config             *services.Config // Store config
	tips               []string         // Array of rotating tips
	currentTipIndex    int              // Current tip index
	tipTimer           time.Time        // Timer for tip rotation
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

	// Initialize tips array
	tips := []string{
		"Tip: Use Tab to load more results!",
		"Tip: Press Esc to go back to search",
		"Tip: Try searching for 'lofi hip hop' or 'cat videos'",
		"Tip: Use ↑/↓ arrows to navigate video list",
		"Tip: Press Enter to play selected video",
	}

	return &SearchComponent{
		textInput:       ti,
		spinner:         s,
		width:           80,
		height:          20,
		cache:           make(map[string][]types.Video),
		recentSearches:  []string{},
		config:          config, // Store config
		tips:            tips,
		currentTipIndex: 0,
		tipTimer:        time.Now(),
	}
}

func (s *SearchComponent) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		s.rotateTip(),
	)
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
	case TipRotateMsg:
		s.currentTipIndex = (s.currentTipIndex + 1) % len(s.tips)
		return s, s.rotateTip()
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *SearchComponent) rotateTip() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return TipRotateMsg{}
	})
}

func (s *SearchComponent) View() string {
	// Logo and branding
	logo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ADD8")).
		Bold(true).
		Render(`
   _____             _                 _______    _          
  / ____|           | |               |__   __|  | |         
 | |  __  ___  _ __ | |__   ___ _ __     | |_   _| |__   ___ 
 | | |_ |/ _ \| '_ \| '_ \ / _ \ '__|    | | | | | '_ \ / _ \
 | |__| | (_) | |_) | | | |  __/ |       | | |_| | |_) |  __/
  \_____|\___/| .__/|_| |_|\___|_|       |_|\__,_|_.__/ \___|
               | |                                            
               |_|                                            `)

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ADD8")).
		Align(lipgloss.Left).
		Width(s.width).
		Render("GopherTube")

	// Tagline
	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Left).
		Width(s.width).
		Render("Fast, private, and keyboard-driven YouTube search & playback")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(s.width).
		Render(strings.Repeat("─", s.width))

	// Search box
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

	// Help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Left).
		Width(s.width).
		Render("Enter: Search  |  Ctrl+C or Esc: Quit")

	// Recent searches (if any)
	recent := ""
	if len(s.recentSearches) > 0 {
		recent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Width(s.width).
			Render("Recent: " + strings.Join(s.recentSearches, ", "))
	}

	// Example queries
	examples := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(s.width).
		Render("Try: lofi hip hop  |  mental outlaw  |  go programming  |  linux tutorial  |  cat videos")

	// Stats
	stats := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(s.width).
		Render(fmt.Sprintf("Cached: %d  |  Last search: %v", len(s.cache), s.lastSearchDuration))

	// Rotating tip
	tip := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ADD8")).
		Italic(true).
		Width(s.width).
		Render(s.tips[s.currentTipIndex])

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		title,
		tagline,
		divider,
		content,
		"",
		help,
		divider,
		recent,
		examples,
		stats,
		tip,
		divider,
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
	return services.SearchYouTube(query, s.config.SearchLimit)
}

func (s *SearchComponent) searchVideos(query string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		// Check cache first
		s.cacheMux.RLock()
		if cached, exists := s.cache[query]; exists {
			s.cacheMux.RUnlock()
			s.lastSearchDuration = time.Since(start)
			s.addRecentSearch(query)
			return SearchResultMsg{Videos: cached}
		}
		s.cacheMux.RUnlock()

		// Use yt-dlp to search for videos
		videos, err := s.searchWithYtDlp(query)
		s.lastSearchDuration = time.Since(start)
		s.addRecentSearch(query)

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

func (s *SearchComponent) addRecentSearch(query string) {
	if query == "" {
		return
	}
	// Only add if not already the most recent
	if len(s.recentSearches) > 0 && s.recentSearches[0] == query {
		return
	}
	// Remove if already present elsewhere
	for i, q := range s.recentSearches {
		if q == query {
			s.recentSearches = append(s.recentSearches[:i], s.recentSearches[i+1:]...)
			break
		}
	}
	// Prepend
	s.recentSearches = append([]string{query}, s.recentSearches...)
	// Limit to 5
	if len(s.recentSearches) > 5 {
		s.recentSearches = s.recentSearches[:5]
	}
}

func (s *SearchComponent) searchWithYtDlp(query string) ([]types.Video, error) {
	// Use the new YouTube scraper instead of yt-dlp
	return services.SearchYouTube(query, s.config.SearchLimit)
}

// ClearCache clears the search cache
func (s *SearchComponent) ClearCache() {
	s.cacheMux.Lock()
	s.cache = make(map[string][]types.Video)
	s.cacheMux.Unlock()
}
