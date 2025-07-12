package components

type Video struct {
	Title       string
	Author      string
	Duration    string
	Views       string
	URL         string
	Thumbnail   string
	Description string
}

type VideoSelectedMsg struct {
	Video Video
}

type VideoPlayedMsg struct {
	Video Video
}

type LoadMoreVideosMsg struct{}

type LoadMoreVideosTimeoutMsg struct{}

type ErrorMsg struct {
	Error string
}
 