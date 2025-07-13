package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// FormatDuration formats seconds into MM:SS or HH:MM:SS
func FormatDuration(seconds int64) string {
	if seconds == 0 {
		return "Unknown"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// FormatViews formats view count into K, M format
func FormatViews(views int64) string {
	if views == 0 {
		return "Unknown"
	}

	if views >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(views)/1000000)
	} else if views >= 1000 {
		return fmt.Sprintf("%.1fK", float64(views)/1000)
	}
	return fmt.Sprintf("%d", views)
}

// TruncateString truncates a string to max length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// CheckCommandExists checks if a command exists in PATH
func CheckCommandExists(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("command '%s' not found in PATH", command)
	}
	return nil
}

// ValidateURL checks if a URL is valid
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// ExtractVideoID extracts YouTube video ID from various URL formats
func ExtractVideoID(url string) string {
	// Common YouTube URL patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`),
		regexp.MustCompile(`youtube\.com/v/([a-zA-Z0-9_-]{11})`),
		regexp.MustCompile(`youtube\.com/watch\?.*v=([a-zA-Z0-9_-]{11})`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// If no pattern matches, assume the URL is already a video ID
	if len(url) == 11 && regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`).MatchString(url) {
		return url
	}

	return ""
}

// SafeGetString safely extracts a string from a map
func SafeGetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// SafeGetInt safely extracts an int64 from a map
func SafeGetInt(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		}
	}
	return 0
}

// Debounce creates a debounced function that delays execution
func Debounce(fn func(), delay time.Duration) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, fn)
	}
}

type PerformanceTimer struct {
	start time.Time
	name  string
}

// StartTimer creates a new performance timer
func StartTimer(name string) *PerformanceTimer {
	return &PerformanceTimer{
		start: time.Now(),
		name:  name,
	}
}

// StopTimer stops the timer and returns the duration
func (pt *PerformanceTimer) StopTimer() time.Duration {
	return time.Since(pt.start)
}

// StopTimerWithLog stops the timer and logs the duration
func (pt *PerformanceTimer) StopTimerWithLog() time.Duration {
	duration := time.Since(pt.start)
	fmt.Printf("[PERF] %s took %v\n", pt.name, duration)
	return duration
}

// FastSplit splits a string by delimiter
func FastSplit(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}

	// Pre-allocate slice
	result := make([]string, 0, len(s)/len(sep)+1)

	start := 0
	for i := 0; i < len(s)-len(sep)+1; i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// FastTrim removes leading and trailing whitespace efficiently
func FastTrim(s string) string {
	start := 0
	end := len(s)

	// Find start of non-whitespace
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Find end of non-whitespace
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
