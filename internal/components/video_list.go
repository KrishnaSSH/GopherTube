package components

import (
	"fmt"
	"time"

	"gophertube/internal/constants"
	"gophertube/internal/types"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VideoList struct {
	videos          []types.Video
	width           int
	height          int
	selected        int
	isLoading       bool
	isPlaying       bool // New state for video playback
	spinner         spinner.Model
	scrollOffset    int      // Track scroll position
	justPlayed      bool     // Track if a video was just played
	tips            []string // Array of rotating tips
	currentTipIndex int      // Current tip index
}

func NewVideoList() *VideoList {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize tips array for video list
	tips := []string{
		"Tip: Press Tab to load more results",
		"Tip: Use ↑/↓ arrows to navigate",
		"Tip: Press Enter to play selected video",
		"Tip: Press Esc to go back to search",
		"Tip: Press 'g' to go to top, 'G' to go to bottom",
	}

	return &VideoList{
		videos:          []types.Video{},
		width:           80,
		height:          20,
		spinner:         s,
		scrollOffset:    0,
		tips:            tips,
		currentTipIndex: 0,
	}
}

func (v *VideoList) Init() tea.Cmd {
	return tea.Batch(
		v.spinner.Tick,
		v.rotateTip(),
	)
}

func (v *VideoList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	availableHeight := v.height - 5
	maxIndex := len(v.videos) - 1

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(v.videos) > 0 && v.selected < len(v.videos) {
				v.isPlaying = true   // Set playing state
				v.justPlayed = false // Reset justPlayed on new play
				return v, tea.Batch(
					v.spinner.Tick,
					func() tea.Msg {
						return VideoSelectedMsg{Video: v.videos[v.selected]}
					},
				)
			}
		case "j", "down":
			if v.selected < maxIndex {
				v.selected++
				// Only scroll if selected is now below the visible window
				if v.selected >= v.scrollOffset+availableHeight {
					v.scrollOffset++
				}
			}
		case "k", "up":
			if v.selected > 0 {
				v.selected--
				// Only scroll if selected is now above the visible window
				if v.selected < v.scrollOffset {
					v.scrollOffset--
				}
			}
		case "g":
			v.selected = 0
			v.scrollOffset = 0
		case "G":
			v.selected = maxIndex
			// Place selected at the bottom of the window if possible
			v.scrollOffset = maxIndex - availableHeight + 1
			if v.scrollOffset < 0 {
				v.scrollOffset = 0
			}
		case "tab":
			if len(v.videos) > 0 && !v.isLoading {
				v.isLoading = true
				return v, tea.Batch(
					v.spinner.Tick,
					func() tea.Msg {
						return LoadMoreVideosMsg{}
					},
					// Add timeout to prevent infinite loading
					tea.Tick(20*time.Second, func(t time.Time) tea.Msg {
						return LoadMoreVideosTimeoutMsg{}
					}),
				)
			}
		}
	case spinner.TickMsg:
		if v.isLoading || v.isPlaying {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			return v, cmd
		}
	case LoadMoreVideosTimeoutMsg:
		if v.isLoading {
			v.isLoading = false
			return v, func() tea.Msg {
				return ErrorMsg{Error: "Loading timeout - try again"}
			}
		}
	case VideoPlayedMsg:
		v.isPlaying = false // Clear playing state
		v.justPlayed = true // Show the "Tab for more" message
		return v, nil
	case TipRotateMsg:
		v.currentTipIndex = (v.currentTipIndex + 1) % len(v.tips)
		return v, v.rotateTip()
	}

	// Clamp selected
	if v.selected < 0 {
		v.selected = 0
	}
	if v.selected > maxIndex {
		v.selected = maxIndex
	}

	// Clamp scrollOffset so the window never goes past the end
	maxScroll := 0
	if availableHeight < len(v.videos) {
		maxScroll = len(v.videos) - availableHeight
	}
	if v.scrollOffset > maxScroll {
		v.scrollOffset = maxScroll
	}
	if v.scrollOffset < 0 {
		v.scrollOffset = 0
	}

	return v, nil
}

func (v *VideoList) rotateTip() tea.Cmd {
	return tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
		return TipRotateMsg{}
	})
}

func (v *VideoList) View() string {
	if len(v.videos) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Left).
			Width(v.width).
			Render(constants.NoVideosMessage)
	}

	availableHeight := v.height - 3 // less padding
	maxScroll := 0
	if availableHeight < len(v.videos) {
		maxScroll = len(v.videos) - availableHeight
	}
	if v.scrollOffset > maxScroll {
		v.scrollOffset = maxScroll
	}
	if v.scrollOffset < 0 {
		v.scrollOffset = 0
	}

	startIndex := v.scrollOffset
	endIndex := startIndex + availableHeight
	if endIndex > len(v.videos) {
		endIndex = len(v.videos)
	}

	var lines []string
	for i := startIndex; i < endIndex; i++ {
		if i < len(v.videos) {
			itemContent := v.renderVideoItem(v.videos[i], i == v.selected)
			if i == v.selected && v.isPlaying {
				itemContent = v.spinner.View() + " " + itemContent
			}
			lines = append(lines, itemContent)
		}
	}

	var content string
	if v.isLoading {
		content = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(v.width).
			Render(v.spinner.View() + " " + constants.LoadingMoreMessage)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Faint(true).
		Align(lipgloss.Left).
		Width(v.width).
		Render(constants.VideoHelpText)

	// Rotating tip
	tip := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ADD8")).
		Faint(false).
		Bold(true).
		Align(lipgloss.Left).
		Width(v.width).
		Render(v.tips[v.currentTipIndex])

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		help,
		tip,
	)
}

func (v *VideoList) renderVideoItem(video types.Video, selected bool) string {
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Faint(true).
		Render(fmt.Sprintf("%s • %s", video.Author, video.Duration))
	title := video.Title
	maxTitleLen := v.width - len(meta) - 4
	if maxTitleLen > 0 && len(title) > maxTitleLen {
		title = title[:maxTitleLen-1] + "…"
	}
	item := title + "  " + meta
	if selected {
		item = lipgloss.NewStyle().
			Background(lipgloss.Color("#00ADD8")).
			Foreground(lipgloss.Color("#222222")).
			MarginRight(10).
			Width(v.width - 10).
			Render(item)
	} else {
		item = lipgloss.NewStyle().Width(v.width).Render(item)
	}
	return item
}

func (v *VideoList) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *VideoList) SetVideos(videos []types.Video) {
	v.videos = videos
	v.selected = 0
	v.scrollOffset = 0
	v.justPlayed = false
}

func (v *VideoList) ResetLoading() {
	v.isLoading = false
}

func (v *VideoList) ResetPlaying() {
	v.isPlaying = false
	v.justPlayed = false
}

func (v *VideoList) AppendVideos(videos []types.Video) {
	v.videos = append(v.videos, videos...)
	v.isLoading = false
	// Don't reset scroll offset when appending - let user continue from where they were
}
