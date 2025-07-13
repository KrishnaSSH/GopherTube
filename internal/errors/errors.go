package errors

import "fmt"

// Error types for better error handling
type (
	PlaybackError struct {
		URL string
		Err error
	}

	YTDlpError struct {
		Command string
		Err     error
	}
)

// Error methods
func (e *PlaybackError) Error() string {
	return fmt.Sprintf("playback failed for URL '%s': %v", e.URL, e.Err)
}

func (e *PlaybackError) Unwrap() error {
	return e.Err
}

func (e *YTDlpError) Error() string {
	return fmt.Sprintf("yt-dlp command '%s' failed: %v", e.Command, e.Err)
}

func (e *YTDlpError) Unwrap() error {
	return e.Err
}

// Helper functions
func NewPlaybackError(url string, err error) *PlaybackError {
	return &PlaybackError{URL: url, Err: err}
}

func NewYTDlpError(command string, err error) *YTDlpError {
	return &YTDlpError{Command: command, Err: err}
}
