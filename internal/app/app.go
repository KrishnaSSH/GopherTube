package app

import (
	"fmt"

	"gophertube/internal/components"
	"gophertube/internal/services"
	"gophertube/internal/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	searchComponent *components.SearchComponent
	videoList       *components.VideoList
	statusBar       *components.StatusBar
	config          *services.Config
	mpvService      *services.MPVService
	state           AppState
}

type AppState struct {
	CurrentView   View
	SearchQuery   string
	Videos        []types.Video
	SelectedIndex int
	IsLoading     bool
	ErrorMessage  string
}

type View int

const (
	ViewSearch View = iota
	ViewVideoList
	ViewVideoPlayer
)

func NewApp() *App {
	config := services.NewConfig()
	mpvService := services.NewMPVService(config)

	return &App{
		searchComponent: components.NewSearchComponent(config),
		videoList:       components.NewVideoList(),
		statusBar:       components.NewStatusBar(),
		config:          config,
		mpvService:      mpvService,
		state: AppState{
			CurrentView:   ViewSearch,
			SelectedIndex: 0,
		},
	}
}

func (a *App) Run() error {
	p := tea.NewProgram(a, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}
	return nil
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.searchComponent.Init(),
		a.videoList.Init(),
		a.statusBar.Init(),
	)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "esc":
			if a.state.CurrentView == ViewVideoList {
				a.state.CurrentView = ViewSearch
				a.state.Videos = nil
				a.state.SelectedIndex = 0
				// Reset search component loading state
				a.searchComponent.ResetLoading()
			} else if a.state.CurrentView == ViewSearch {
				return a, tea.Quit
			}
		}

	case components.SearchResultMsg:
		if msg.Error != "" {
			a.state.ErrorMessage = msg.Error
			a.state.IsLoading = false
			a.statusBar.SetMessage("Search failed: " + msg.Error)
			return a, nil
		}

		// Check if this is a "load more" operation
		if a.state.CurrentView == ViewVideoList && len(a.state.Videos) > 0 {
			// Append new videos to existing list
			a.videoList.AppendVideos(msg.Videos)
			a.state.Videos = append(a.state.Videos, msg.Videos...)
			a.statusBar.SetMessage(fmt.Sprintf("Added %d more videos", len(msg.Videos)))
		} else {
			// Replace videos for new search
			a.state.Videos = msg.Videos
			a.state.CurrentView = ViewVideoList
			a.state.IsLoading = false
			a.state.ErrorMessage = ""
			a.videoList.SetVideos(msg.Videos)
			a.statusBar.SetMessage(fmt.Sprintf("Found %d videos", len(msg.Videos)))
		}
		return a, nil

	case components.VideoSelectedMsg:
		return a, a.playVideo(msg.Video)

	case components.VideoPlayedMsg:
		// Video playback completed, clear playing state
		return a, nil

	case components.ErrorMsg:
		a.state.ErrorMessage = msg.Error
		a.state.IsLoading = false
		// Reset loading state in video list if it's loading
		a.videoList.ResetLoading()
		a.videoList.ResetPlaying() // Also reset playing state
		a.statusBar.SetMessage("Error: " + msg.Error)
		return a, nil

	case components.LoadMoreVideosMsg:
		return a, a.loadMoreVideos()

	case tea.WindowSizeMsg:
		a.updateLayout(msg.Width, msg.Height)
		return a, nil
	}

	// Update current view component
	switch a.state.CurrentView {
	case ViewSearch:
		searchModel, cmd := a.searchComponent.Update(msg)
		a.searchComponent = searchModel.(*components.SearchComponent)
		return a, cmd
	case ViewVideoList:
		// Update video list for all messages, not just key events
		videoListModel, cmd := a.videoList.Update(msg)
		a.videoList = videoListModel.(*components.VideoList)
		return a, cmd
	}

	return a, nil
}

func (a *App) View() string {
	if a.state.ErrorMessage != "" {
		return lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(60).
			Padding(2, 4).
			Render("Error: " + a.state.ErrorMessage)
	}

	var view string
	switch a.state.CurrentView {
	case ViewSearch:
		view = a.searchComponent.View()
	case ViewVideoList:
		view = a.videoList.View()
	}

	// Only add status bar if it has content and we're not in video list view
	statusBar := a.statusBar.View()
	if statusBar != "" && a.state.CurrentView != ViewVideoList {
		return lipgloss.NewStyle().Padding(2, 4).Render(lipgloss.JoinVertical(lipgloss.Left, view, statusBar))
	}
	return lipgloss.NewStyle().Padding(2, 4).Render(view)
}

func (a *App) updateLayout(width, height int) {
	a.searchComponent.SetSize(width, height-3) // Reserve space for status bar
	a.videoList.SetSize(width, height-3)
	a.statusBar.SetSize(width, 3)
}

func (a *App) playVideo(video types.Video) tea.Cmd {
	return func() tea.Msg {
		if err := a.mpvService.PlayVideo(video.URL); err != nil {
			return components.ErrorMsg{Error: err.Error()}
		}
		return components.VideoPlayedMsg{Video: video}
	}
}

func (a *App) loadMoreVideos() tea.Cmd {
	return func() tea.Msg {
		// Use the current search query to get more videos
		query := a.searchComponent.GetCurrentQuery()
		if query == "" {
			return components.ErrorMsg{Error: "No search query available"}
		}

		// Get more videos using yt-dlp
		videos, err := a.searchComponent.SearchWithQuery(query)
		if err != nil {
			return components.ErrorMsg{Error: err.Error()}
		}

		return components.SearchResultMsg{Videos: videos}
	}
}
