package components

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	width     int
	height    int
	message   string
	timestamp time.Time
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		width:     80,
		height:    3,
		message:   "Ready",
		timestamp: time.Now(),
	}
}

func (s *StatusBar) Init() tea.Cmd {
	return nil
}

func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s.message = fmt.Sprintf("Key pressed: %s", msg.String())
		s.timestamp = time.Now()
	}

	return s, nil
}

func (s *StatusBar) View() string {
	if s.message == "Ready" || s.message == "" {
		return ""
	}
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	timeStr := s.timestamp.Format("15:04:05")
	status := fmt.Sprintf("%s | %s", s.message, timeStr)

	return statusStyle.Width(s.width).Render(status)
}

func (s *StatusBar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *StatusBar) SetMessage(message string) {
	s.message = message
	s.timestamp = time.Now()
}
