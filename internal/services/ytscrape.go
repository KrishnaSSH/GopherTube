package services

import (
	"bufio"
	"context"
	"errors"
	"gophertube/internal/types"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

var ytInitialDataPrefix = "var ytInitialData = "

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func SearchYouTube(query string, limit int) ([]types.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YouTube returned non-200 status")
	}

	jsonData, err := extractYtInitialData(resp.Body)
	if err != nil {
		return nil, err
	}

	results := gjson.GetBytes(jsonData, "..videoRenderer")
	if !results.Exists() {
		return nil, errors.New("no videos found in YouTube search")
	}

	videos := make([]types.Video, 0, limit)
	ch := make(chan types.Video, limit)
	var collected int32
	workerCount := runtime.NumCPU()
	sem := make(chan struct{}, workerCount)
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	for _, vr := range results.Array() {
		if atomic.LoadInt32(&collected) >= int32(limit) {
			break
		}
		sem <- struct{}{}
		go func(vr gjson.Result) {
			defer func() { <-sem }()
			if ctx2.Err() != nil {
				return
			}
			video := parseVideoRendererGJSON(vr)
			if video.Title != "" && video.URL != "" {
				if atomic.AddInt32(&collected, 1) <= int32(limit) {
					ch <- video
					if atomic.LoadInt32(&collected) == int32(limit) {
						cancel2()
					}
				}
			}
		}(vr)
	}

outer:
	for range results.Array() {
		if len(videos) >= limit {
			break
		}
		select {
		case v := <-ch:
			videos = append(videos, v)
		case <-ctx2.Done():
			break outer
		}
	}

	if len(videos) == 0 {
		return nil, errors.New("no videos found in YouTube search")
	}
	return videos, nil
}

// Fast line-by-line extraction of ytInitialData JSON
func extractYtInitialData(r io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, ytInitialDataPrefix); idx != -1 {
			jsonStart := idx + len(ytInitialDataPrefix)
			jsonLine := line[jsonStart:]
			// If the JSON is not complete on this line, keep reading
			for !strings.HasSuffix(strings.TrimSpace(jsonLine), ";") && scanner.Scan() {
				jsonLine += scanner.Text()
			}
			jsonLine = strings.TrimSuffix(jsonLine, ";")
			return []byte(jsonLine), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, errors.New("ytInitialData not found in YouTube page")
}

func urlQueryEscape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "+"), "#", "%23")
}

func parseVideoRendererGJSON(vr gjson.Result) types.Video {
	title := vr.Get("title.runs.0.text").String()
	videoId := vr.Get("videoId").String()
	url := "https://www.youtube.com/watch?v=" + videoId
	channel := vr.Get("longBylineText.runs.0.text").String()
	duration := vr.Get("lengthText.simpleText").String()
	thumb := ""
	thumbs := vr.Get("thumbnail.thumbnails")
	if thumbs.Exists() && thumbs.IsArray() && len(thumbs.Array()) > 0 {
		thumb = thumbs.Array()[len(thumbs.Array())-1].Get("url").String()
	}
	return types.Video{
		Title:     title,
		URL:       url,
		Author:    channel,
		Duration:  duration,
		Thumbnail: thumb,
	}
}
