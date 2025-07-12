package utils

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ResponsiveLayout handles dynamic sizing for different terminal dimensions
type ResponsiveLayout struct {
	Width  int
	Height int
}

// NewResponsiveLayout creates a new responsive layout manager
func NewResponsiveLayout(width, height int) *ResponsiveLayout {
	return &ResponsiveLayout{
		Width:  width,
		Height: height,
	}
}

// GetResponsiveWidth returns a width that adapts to terminal size
func (r *ResponsiveLayout) GetResponsiveWidth() int {
	// Ensure minimum width for usability
	minWidth := 40
	maxWidth := r.Width - 4 // Leave some margin

	if maxWidth < minWidth {
		return minWidth
	}

	return maxWidth
}

// GetResponsiveHeight returns a height that adapts to terminal size
func (r *ResponsiveLayout) GetResponsiveHeight() int {
	// Ensure minimum height for usability
	minHeight := 10
	maxHeight := r.Height - 4 // Leave some margin

	if maxHeight < minHeight {
		return minHeight
	}

	return maxHeight
}

// GetContentWidth returns the width available for content
func (r *ResponsiveLayout) GetContentWidth() int {
	return r.GetResponsiveWidth() - 4 // Account for padding
}

// GetContentHeight returns the height available for content
func (r *ResponsiveLayout) GetContentHeight() int {
	return r.GetResponsiveHeight() - 6 // Account for title, help text, and spacing
}

// IsSmallTerminal returns true if terminal is small
func (r *ResponsiveLayout) IsSmallTerminal() bool {
	return r.Width < 60 || r.Height < 15
}

// IsMediumTerminal returns true if terminal is medium sized
func (r *ResponsiveLayout) IsMediumTerminal() bool {
	return (r.Width >= 60 && r.Width < 100) || (r.Height >= 15 && r.Height < 25)
}

// IsLargeTerminal returns true if terminal is large
func (r *ResponsiveLayout) IsLargeTerminal() bool {
	return r.Width >= 100 && r.Height >= 25
}

// GetOptimalPadding returns padding based on terminal size
func (r *ResponsiveLayout) GetOptimalPadding() (int, int) {
	if r.IsSmallTerminal() {
		return 1, 2
	} else if r.IsMediumTerminal() {
		return 2, 3
	} else {
		return 3, 4
	}
}

// GetOptimalMargin returns margin based on terminal size
func (r *ResponsiveLayout) GetOptimalMargin() int {
	if r.IsSmallTerminal() {
		return 1
	} else if r.IsMediumTerminal() {
		return 2
	} else {
		return 3
	}
}

// TruncateText truncates text to fit within the given width
func (r *ResponsiveLayout) TruncateText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	// Leave space for ellipsis
	if maxWidth <= 3 {
		return "..."
	}

	// For very small terminals, be more aggressive with truncation
	if r.IsSmallTerminal() && maxWidth < 20 {
		return text[:maxWidth-3] + "..."
	}

	// For medium terminals, try to preserve more text
	if r.IsMediumTerminal() && maxWidth < 40 {
		return text[:maxWidth-3] + "..."
	}

	return text[:maxWidth-3] + "..."
}

// WrapText wraps text to fit within the given width
func (r *ResponsiveLayout) WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxWidth {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// GetResponsiveStyle returns a style that adapts to terminal size
func (r *ResponsiveLayout) GetResponsiveStyle() lipgloss.Style {
	paddingTop, paddingLeft := r.GetOptimalPadding()
	margin := r.GetOptimalMargin()

	return lipgloss.NewStyle().
		Padding(paddingTop, paddingLeft).
		Margin(margin)
}

// GetTitleStyle returns a responsive title style
func (r *ResponsiveLayout) GetTitleStyle() lipgloss.Style {
	width := r.GetResponsiveWidth()

	style := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Left).
		Width(width)

	if r.IsLargeTerminal() {
		style = style.Padding(1, 0)
	}

	return style
}

// GetContentStyle returns a responsive content style
func (r *ResponsiveLayout) GetContentStyle() lipgloss.Style {
	width := r.GetResponsiveWidth()

	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(width)
}

// GetHelpStyle returns a responsive help text style
func (r *ResponsiveLayout) GetHelpStyle() lipgloss.Style {
	width := r.GetResponsiveWidth()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Left).
		Width(width)
}

// GetVideoItemStyle returns a responsive video item style
func (r *ResponsiveLayout) GetVideoItemStyle(selected bool) lipgloss.Style {
	width := r.GetContentWidth()

	style := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left)

	if selected {
		style = style.Background(lipgloss.Color("#444444"))
	}

	if r.IsSmallTerminal() {
		style = style.Padding(0, 1)
	} else {
		style = style.Padding(0, 2)
	}

	return style
}

// GetStatusStyle returns a responsive status bar style
func (r *ResponsiveLayout) GetStatusStyle() lipgloss.Style {
	width := r.GetResponsiveWidth()

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(width).
		Align(lipgloss.Left)

	if r.IsSmallTerminal() {
		style = style.Padding(0, 1)
	} else {
		style = style.Padding(0, 2)
	}

	return style
}

// CalculateOptimalColumns calculates the optimal number of columns for content
func (r *ResponsiveLayout) CalculateOptimalColumns() int {
	width := r.GetContentWidth()

	if width < 60 {
		return 1
	} else if width < 100 {
		return 2
	} else {
		return 3
	}
}

// GetColumnWidth calculates the width for each column
func (r *ResponsiveLayout) GetColumnWidth() int {
	columns := r.CalculateOptimalColumns()
	contentWidth := r.GetContentWidth()

	return int(math.Floor(float64(contentWidth) / float64(columns)))
}

// IsCompactMode returns true if we should use compact mode
func (r *ResponsiveLayout) IsCompactMode() bool {
	return r.IsSmallTerminal() || (r.Width < 80 && r.Height < 20)
}

// GetCompactPadding returns padding for compact mode
func (r *ResponsiveLayout) GetCompactPadding() (int, int) {
	if r.IsCompactMode() {
		return 0, 1
	}
	return r.GetOptimalPadding()
}

// GetOptimalSpacing returns spacing based on terminal size
func (r *ResponsiveLayout) GetOptimalSpacing() int {
	if r.IsSmallTerminal() {
		return 1
	} else if r.IsMediumTerminal() {
		return 2
	} else {
		return 3
	}
}

// GetResponsiveTitle returns a title that adapts to terminal size
func (r *ResponsiveLayout) GetResponsiveTitle(title string) string {
	if r.IsSmallTerminal() {
		return r.TruncateText(title, r.GetContentWidth())
	}
	return title
}
