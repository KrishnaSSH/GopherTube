package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gophertube/internal/types"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

var (
	ytInitialDataPrefix = "var ytInitialData = "
	ytInitialDataSuffix = ";"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func SearchYouTube(query string, limit int) ([]types.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Use a modern Chrome UA
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("YouTube request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube returned non-200 status: %d", resp.StatusCode)
	}

	jsonData, htmlDebug, err := robustExtractYtInitialData(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ytInitialData extraction failed: %w\nFirst 2KB of HTML: %s", err, htmlDebug)
	}

	results := gjson.GetBytes(jsonData, "..videoRenderer")
	if !results.Exists() {
		return nil, errors.New("no videos found in YouTube search (ytInitialData parsed, but no videoRenderer nodes found)")
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

	for range results.Array() {
		if len(videos) >= limit {
			break
		}
		select {
		case v := <-ch:
			videos = append(videos, v)
		case <-ctx2.Done():
			break
		}
	}

	if len(videos) == 0 {
		return nil, errors.New("no videos found in YouTube search (videoRenderer nodes exist, but none parsed)")
	}
	return videos, nil
}

// Tries prefix search, then fallback to searching for ytInitialData anywhere in the HTML
// Returns: (jsonData, htmlDebug, error)
func robustExtractYtInitialData(r io.Reader) ([]byte, string, error) {
	const chunkSize = 32 * 1024
	const debugLen = 2048
	var buf bytes.Buffer
	var htmlDebug string
	found := false
	prefixIdx := -1
	for {
		chunk := make([]byte, chunkSize)
		n, err := r.Read(chunk)
		if n > 0 {
			buf.Write(chunk[:n])
		}
		if buf.Len() >= debugLen && htmlDebug == "" {
			htmlDebug = string(buf.Bytes()[:debugLen])
		}
		if !found {
			b := buf.Bytes()
			prefixIdx = bytes.Index(b, []byte(ytInitialDataPrefix))
			if prefixIdx != -1 {
				buf.Next(prefixIdx + len(ytInitialDataPrefix))
				found = true
				break
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, htmlDebug, err
		}
	}
	if found {
		// Now read until suffix
		var jsonBuf bytes.Buffer
		for {
			chunk := make([]byte, chunkSize)
			n, err := r.Read(chunk)
			if n > 0 {
				if idx := bytes.Index(chunk[:n], []byte(ytInitialDataSuffix)); idx != -1 {
					jsonBuf.Write(chunk[:idx])
					break
				}
				jsonBuf.Write(chunk[:n])
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, htmlDebug, err
			}
		}
		return jsonBuf.Bytes(), htmlDebug, nil
	}
	// Fallback: search for ytInitialData anywhere in the HTML
	b := buf.Bytes()
	idx := bytes.Index(b, []byte("ytInitialData"))
	if idx == -1 {
		return nil, htmlDebug, errors.New("ytInitialData not found in HTML")
	}
	// Try to extract the nearest JSON object after ytInitialData
	start := bytes.Index(b[idx:], []byte("{"))
	end := bytes.Index(b[idx:], []byte("};"))
	if start == -1 || end == -1 || end <= start {
		return nil, htmlDebug, errors.New("ytInitialData JSON object not found after fallback")
	}
	return b[idx+start : idx+end], htmlDebug, nil
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
