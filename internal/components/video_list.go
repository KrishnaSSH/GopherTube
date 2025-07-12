package components

import (
	"fmt"
	"time"

	"gophertube/internal/constants"
	"gophertube/internal/types"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VideoList struct {
	list         list.Model
	videos       []types.Video
	width        int
	height       int
	selected     int
	isLoading    bool
	spinner      spinner.Model
	scrollOffset int // Track scroll position
}

func NewVideoList() *VideoList {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &VideoList{
		list:         l,
		videos:       []types.Video{},
		width:        80,
		height:       20,
		spinner:      s,
		scrollOffset: 0,
	}
}

func (v *VideoList) Init() tea.Cmd {
	return nil
}

func (v *VideoList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(v.videos) > 0 && v.selected < len(v.videos) {
				return v, func() tea.Msg {
					return VideoSelectedMsg{Video: v.videos[v.selected]}
				}
			}
		case "j", "down":
			if v.selected < len(v.videos)-1 {
				v.selected++
				// Auto-scroll if selection goes below visible area
				visibleHeight := v.height - 5 // Account for title, spacing, and help
				if v.selected >= v.scrollOffset+visibleHeight {
					v.scrollOffset = v.selected - visibleHeight + 1
				}
			}
		case "k", "up":
			if v.selected > 0 {
				v.selected--
				// Auto-scroll if selection goes above visible area
				if v.selected < v.scrollOffset {
					v.scrollOffset = v.selected
				}
			}
		case "g":
			v.selected = 0
			v.scrollOffset = 0
		case "G":
			v.selected = len(v.videos) - 1
			// Scroll to show the last item
			visibleHeight := v.height - 5
			if v.selected >= visibleHeight {
				v.scrollOffset = v.selected - visibleHeight + 1
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
		if v.isLoading {
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
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v *VideoList) View() string {
	if len(v.videos) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Left).
			Width(v.width).
			Render(constants.NoVideosMessage)
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Left).
		Width(v.width).
		Render("Results")

	var videoItems []string
	for i, video := range v.videos {
		videoItems = append(videoItems, v.renderVideoItem(video, i == v.selected))
	}

	// Calculate visible area
	visibleHeight := v.height - 5 // Account for title, spacing, and help

	// Apply scroll offset and limit to visible height
	startIndex := v.scrollOffset
	endIndex := startIndex + visibleHeight
	if endIndex > len(videoItems) {
		endIndex = len(videoItems)
	}

	// Ensure we don't go out of bounds
	if startIndex >= len(videoItems) {
		startIndex = 0
		v.scrollOffset = 0
	}

	// Get the visible portion of the list
	if startIndex < len(videoItems) {
		videoItems = videoItems[startIndex:endIndex]
	} else {
		videoItems = []string{}
	}

	// Add scroll indicators if needed
	if v.scrollOffset > 0 {
		// Show indicator that there are items above
		scrollUpIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(v.width).
			Render("↑ More videos above")
		videoItems = append([]string{scrollUpIndicator}, videoItems...)
	}

	if v.scrollOffset+visibleHeight < len(v.videos) {
		// Show indicator that there are items below
		scrollDownIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(v.width).
			Render("↓ More videos below")
		videoItems = append(videoItems, scrollDownIndicator)
	}

	var content string
	if v.isLoading {
		content = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(v.width).
			Render(v.spinner.View() + " " + constants.LoadingMoreMessage)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left, videoItems...)
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Left).
		Width(v.width).
		Render(constants.VideoHelpText)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		help,
	)
}

func (v *VideoList) renderVideoItem(video types.Video, selected bool) string {
	// Title (bold if selected)
	titleStyle := lipgloss.NewStyle()
	if selected {
		titleStyle = titleStyle.Bold(true).Underline(true)
	}
	title := titleStyle.Render(lipgloss.NewStyle().MaxWidth(v.width - 2).Render(video.Title))

	// Author and duration (dimmed)
	authorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Faint(true)
	author := authorStyle.Render(fmt.Sprintf("%s • %s", video.Author, video.Duration))

	return lipgloss.JoinVertical(lipgloss.Left, title, author)
}

func (v *VideoList) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *VideoList) SetVideos(videos []types.Video) {
	v.videos = videos
	v.selected = 0
	v.scrollOffset = 0
}

func (v *VideoList) ResetLoading() {
	v.isLoading = false
}

func (v *VideoList) AppendVideos(videos []types.Video) {
	v.videos = append(v.videos, videos...)
	v.isLoading = false
	// Don't reset scroll offset when appending - let user continue from where they were
}
