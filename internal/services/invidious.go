package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"gophertube/internal/errors"
	"gophertube/internal/types"
	"gophertube/internal/utils"
)

type InvidiousService struct {
	config *Config
}

type InvidiousVideo struct {
	VideoID       string `json:"videoId"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	AuthorID      string `json:"authorId"`
	LengthSeconds int64  `json:"lengthSeconds"`
	ViewCount     int64  `json:"viewCount"`
	Published     int64  `json:"published"`
	Description   string `json:"description"`
	Thumbnail     string `json:"videoThumbnails"`
}

type InvidiousSearchResponse struct {
	Videos []InvidiousVideo `json:"content"`
}

func NewInvidiousService(config *Config) *InvidiousService {
	return &InvidiousService{
		config: config,
	}
}

func (i *InvidiousService) SearchVideos(query string) ([]types.Video, error) {
	// Performance timer
	timer := utils.StartTimer("invidious_search")
	defer timer.StopTimerWithLog()

	// Check if curl and jq are available
	if err := utils.CheckCommandExists("curl"); err != nil {
		return nil, errors.NewSearchError(query, fmt.Errorf("curl not found: %w", err))
	}
	if err := utils.CheckCommandExists("jq"); err != nil {
		return nil, errors.NewSearchError(query, fmt.Errorf("jq not found: %w", err))
	}

	// Build the search URL
	searchURL := fmt.Sprintf("%s/api/v1/search?q=%s&type=video&sort_by=relevance&date=all&duration=all&features=all&region=all&limit=%d",
		strings.TrimSuffix(i.config.InvidiousURL, "/"),
		url.QueryEscape(query),
		i.config.SearchLimit,
	)

	// Create curl command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "curl",
		"-s",               // Silent mode
		"-L",               // Follow redirects
		"--max-time", "10", // 10 second timeout
		"--connect-timeout", "5", // 5 second connection timeout
		"-H", "User-Agent: Mozilla/5.0 (compatible; GopherTube/1.0)",
		searchURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.NewSearchError(query, fmt.Errorf("curl failed: %w", err))
	}

	// Parse JSON response using jq for better performance
	return i.parseSearchResults(string(output), query)
}

func (i *InvidiousService) parseSearchResults(jsonData string, query string) ([]types.Video, error) {
	// Use jq to extract and format the data efficiently
	jqCmd := exec.Command("jq", "-r",
		`.content[] | {
			"videoId": .videoId,
			"title": .title,
			"author": .author,
			"lengthSeconds": .lengthSeconds,
			"viewCount": .viewCount,
			"description": .description
		}`,
	)

	jqCmd.Stdin = strings.NewReader(jsonData)
	output, err := jqCmd.Output()
	if err != nil {
		// Fallback to Go JSON parsing if jq fails
		return i.parseSearchResultsGo(jsonData, query)
	}

	// Parse the jq output
	var videos []types.Video
	lines := utils.FastSplit(utils.FastTrim(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var videoData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &videoData); err != nil {
			continue
		}

		video := types.Video{
			Title:       utils.SafeGetString(videoData, "title"),
			Author:      utils.SafeGetString(videoData, "author"),
			Duration:    utils.FormatDuration(utils.SafeGetInt(videoData, "lengthSeconds")),
			Views:       utils.FormatViews(utils.SafeGetInt(videoData, "viewCount")),
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", utils.SafeGetString(videoData, "videoId")),
			Thumbnail:   fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", utils.SafeGetString(videoData, "videoId")),
			Description: utils.SafeGetString(videoData, "description"),
		}

		// Only add videos with valid URLs and titles
		if video.URL != "" && video.Title != "" {
			videos = append(videos, video)
		}
	}

	if len(videos) == 0 {
		return nil, errors.NewSearchError(query, fmt.Errorf("no videos found"))
	}

	return videos, nil
}

func (i *InvidiousService) parseSearchResultsGo(jsonData string, query string) ([]types.Video, error) {
	var response InvidiousSearchResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		return nil, errors.NewSearchError(query, fmt.Errorf("failed to parse JSON: %w", err))
	}

	var videos []types.Video
	for _, invVideo := range response.Videos {
		video := types.Video{
			Title:       invVideo.Title,
			Author:      invVideo.Author,
			Duration:    utils.FormatDuration(invVideo.LengthSeconds),
			Views:       utils.FormatViews(invVideo.ViewCount),
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", invVideo.VideoID),
			Thumbnail:   fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", invVideo.VideoID),
			Description: invVideo.Description,
		}

		// Only add videos with valid URLs and titles
		if video.URL != "" && video.Title != "" {
			videos = append(videos, video)
		}
	}

	if len(videos) == 0 {
		return nil, errors.NewSearchError(query, fmt.Errorf("no videos found"))
	}

	return videos, nil
}

// GetVideoInfo gets detailed information about a video using Invidious
func (i *InvidiousService) GetVideoInfo(videoID string) (*types.Video, error) {
	// Extract video ID from URL if needed
	if strings.Contains(videoID, "youtube.com") {
		videoID = utils.ExtractVideoID(videoID)
	}

	// Build the video info URL
	videoURL := fmt.Sprintf("%s/api/v1/videos/%s",
		strings.TrimSuffix(i.config.InvidiousURL, "/"),
		videoID,
	)

	// Create curl command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "curl",
		"-s",
		"-L",
		"--max-time", "8",
		"--connect-timeout", "5",
		"-H", "User-Agent: Mozilla/5.0 (compatible; GopherTube/1.0)",
		videoURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.NewSearchError(videoID, fmt.Errorf("failed to get video info: %w", err))
	}

	// Parse the response
	var videoData map[string]interface{}
	if err := json.Unmarshal(output, &videoData); err != nil {
		return nil, errors.NewSearchError(videoID, fmt.Errorf("failed to parse video info: %w", err))
	}

	video := &types.Video{
		Title:       utils.SafeGetString(videoData, "title"),
		Author:      utils.SafeGetString(videoData, "author"),
		Duration:    utils.FormatDuration(utils.SafeGetInt(videoData, "lengthSeconds")),
		Views:       utils.FormatViews(utils.SafeGetInt(videoData, "viewCount")),
		URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Thumbnail:   fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", videoID),
		Description: utils.SafeGetString(videoData, "description"),
	}

	return video, nil
}
