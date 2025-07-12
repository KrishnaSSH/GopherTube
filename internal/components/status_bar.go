package components

import (
	"fmt"
	"gophertube/internal/utils"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StatusBar struct {
	layout    *utils.ResponsiveLayout
	message   string
	timestamp time.Time
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		layout:    utils.NewResponsiveLayout(80, 3),
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
	statusStyle := s.layout.GetStatusStyle()

	timeStr := s.timestamp.Format("15:04:05")
	status := fmt.Sprintf("%s | %s", s.message, timeStr)

	return statusStyle.Render(status)
}

func (s *StatusBar) SetSize(width, height int) {
	s.layout = utils.NewResponsiveLayout(width, height)
}

func (s *StatusBar) SetMessage(message string) {
	s.message = message
	s.timestamp = time.Now()
}
