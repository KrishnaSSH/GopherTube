package utils

import (
	"testing"
)

func TestResponsiveLayout(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected struct {
			isSmall  bool
			isMedium bool
			isLarge  bool
		}
	}{
		{
			name:   "Small terminal",
			width:  50,
			height: 10,
			expected: struct {
				isSmall  bool
				isMedium bool
				isLarge  bool
			}{
				isSmall:  true,
				isMedium: false,
				isLarge:  false,
			},
		},
		{
			name:   "Medium terminal",
			width:  80,
			height: 20,
			expected: struct {
				isSmall  bool
				isMedium bool
				isLarge  bool
			}{
				isSmall:  false,
				isMedium: true,
				isLarge:  false,
			},
		},
		{
			name:   "Large terminal",
			width:  120,
			height: 30,
			expected: struct {
				isSmall  bool
				isMedium bool
				isLarge  bool
			}{
				isSmall:  false,
				isMedium: false,
				isLarge:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewResponsiveLayout(tt.width, tt.height)

			if layout.IsSmallTerminal() != tt.expected.isSmall {
				t.Errorf("IsSmallTerminal() = %v, want %v", layout.IsSmallTerminal(), tt.expected.isSmall)
			}

			if layout.IsMediumTerminal() != tt.expected.isMedium {
				t.Errorf("IsMediumTerminal() = %v, want %v", layout.IsMediumTerminal(), tt.expected.isMedium)
			}

			if layout.IsLargeTerminal() != tt.expected.isLarge {
				t.Errorf("IsLargeTerminal() = %v, want %v", layout.IsLargeTerminal(), tt.expected.isLarge)
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	layout := NewResponsiveLayout(80, 20)

	tests := []struct {
		name     string
		text     string
		maxWidth int
		expected string
	}{
		{
			name:     "Short text",
			text:     "Hello",
			maxWidth: 10,
			expected: "Hello",
		},
		{
			name:     "Long text",
			text:     "This is a very long text that needs to be truncated",
			maxWidth: 20,
			expected: "This is a very lo...",
		},
		{
			name:     "Very short max width",
			text:     "Hello",
			maxWidth: 3,
			expected: "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := layout.TruncateText(tt.text, tt.maxWidth)
			if result != tt.expected {
				t.Errorf("TruncateText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetResponsiveWidth(t *testing.T) {
	layout := NewResponsiveLayout(100, 20)
	width := layout.GetResponsiveWidth()

	// Should be less than or equal to terminal width
	if width > 100 {
		t.Errorf("GetResponsiveWidth() = %v, should be <= 100", width)
	}

	// Should be at least minimum width
	if width < 40 {
		t.Errorf("GetResponsiveWidth() = %v, should be >= 40", width)
	}
}

func TestGetResponsiveHeight(t *testing.T) {
	layout := NewResponsiveLayout(80, 25)
	height := layout.GetResponsiveHeight()

	// Should be less than or equal to terminal height
	if height > 25 {
		t.Errorf("GetResponsiveHeight() = %v, should be <= 25", height)
	}

	// Should be at least minimum height
	if height < 10 {
		t.Errorf("GetResponsiveHeight() = %v, should be >= 10", height)
	}
}
