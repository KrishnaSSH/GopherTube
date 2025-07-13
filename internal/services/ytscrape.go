package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophertube/internal/types"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// SearchYouTube scrapes YouTube search results page and returns a list of videos
func SearchYouTube(query string, limit int) ([]types.Video, error) {
	url := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch YouTube: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read YouTube response: %w", err)
	}

	// Extract ytInitialData JSON
	re := regexp.MustCompile(`(?s)var ytInitialData = (\{.*?\});`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return nil, errors.New("ytInitialData not found in YouTube page")
	}
	jsonData := matches[1]

	// Parse JSON
	var root map[string]interface{}
	if err := json.Unmarshal(jsonData, &root); err != nil {
		return nil, fmt.Errorf("failed to parse ytInitialData: %w", err)
	}

	// Traverse to videoRenderer nodes
	videos := []types.Video{}
	var walk func(interface{})
	walk = func(node interface{}) {
		if len(videos) >= limit {
			return
		}
		m, ok := node.(map[string]interface{})
		if ok {
			if vr, ok := m["videoRenderer"]; ok {
				video := parseVideoRenderer(vr)
				if video.Title != "" && video.URL != "" {
					videos = append(videos, video)
				}
			}
			for _, v := range m {
				walk(v)
			}
		} else if arr, ok := node.([]interface{}); ok {
			for _, v := range arr {
				walk(v)
			}
		}
	}
	walk(root)

	if len(videos) == 0 {
		return nil, errors.New("no videos found in YouTube search")
	}
	if len(videos) > limit {
		videos = videos[:limit]
	}
	return videos, nil
}

func urlQueryEscape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "+"), "#", "%23")
}

func parseVideoRenderer(vr interface{}) types.Video {
	m, ok := vr.(map[string]interface{})
	if !ok {
		return types.Video{}
	}
	title := safeJQText(m, "title", "runs", 0, "text")
	videoId := safeJQString(m, "videoId")
	url := "https://www.youtube.com/watch?v=" + videoId
	channel := safeJQText(m, "longBylineText", "runs", 0, "text")
	duration := safeJQString(m, "lengthText", "simpleText")
	thumb := ""
	if thumbs, ok := m["thumbnail"].(map[string]interface{}); ok {
		if arr, ok := thumbs["thumbnails"].([]interface{}); ok && len(arr) > 0 {
			if t, ok := arr[len(arr)-1].(map[string]interface{}); ok {
				thumb, _ = t["url"].(string)
			}
		}
	}
	return types.Video{
		Title:     title,
		URL:       url,
		Author:    channel,
		Duration:  duration,
		Thumbnail: thumb,
	}
}

func safeJQString(m map[string]interface{}, keys ...string) string {
	cur := m
	for i, k := range keys {
		if i == len(keys)-1 {
			if v, ok := cur[k].(string); ok {
				return v
			}
			return ""
		}
		if v, ok := cur[k].(map[string]interface{}); ok {
			cur = v
		} else {
			return ""
		}
	}
	return ""
}

func safeJQText(m map[string]interface{}, k1, k2 string, idx int, k3 string) string {
	if a, ok := m[k1].(map[string]interface{}); ok {
		if arr, ok := a[k2].([]interface{}); ok && len(arr) > idx {
			if t, ok := arr[idx].(map[string]interface{}); ok {
				if s, ok := t[k3].(string); ok {
					return s
				}
			}
		}
	}
	return ""
}
