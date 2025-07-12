package components

import "gophertube/internal/types"

// Message types for component communication
type VideoSelectedMsg struct {
	Video types.Video
}

type VideoPlayedMsg struct {
	Video types.Video
}

type ErrorMsg struct {
	Error string
}

type LoadMoreVideosMsg struct{}

type LoadMoreVideosTimeoutMsg struct{}
