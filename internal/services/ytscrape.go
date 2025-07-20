package services

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"gophertube/internal/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Create a custom HTTP client with optimizations for faster downloads
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		ForceAttemptHTTP2:   true,
		MaxConnsPerHost:     0,
		DisableKeepAlives:   false,
	},
}

func SearchYouTube(query string, limit int, progress func(current, total int)) ([]types.Video, error) {
	if progress != nil {
		progress(0, 1)
	}

	url := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query) + "&sp=EgIQAQ%253D%253D"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if progress != nil {
		progress(1, 2)
	}

	re := regexp.MustCompile(`(?s)var ytInitialData = (\{.*?\});`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return nil, errors.New("ytInitialData not found")
	}
	var root map[string]interface{}
	if err := json.Unmarshal(m[1], &root); err != nil {
		return nil, err
	}

	videos := []types.Video{}
	var walk func(interface{})
	walk = func(node interface{}) {
		if m, ok := node.(map[string]interface{}); ok {
			if vr, ok := m["videoRenderer"]; ok {
				v := parseVideoRenderer(vr)
				if v.Title != "" && v.URL != "" {
					videos = append(videos, v)
				}
			}
			if cr, ok := m["compactVideoRenderer"]; ok {
				v := parseVideoRenderer(cr)
				if v.Title != "" && v.URL != "" {
					videos = append(videos, v)
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
		return nil, errors.New("no videos found")
	}

	if len(videos) < limit {
		searchStrategies := []string{
			"",
			"&sp=EgIQAQ%253D%253D",
			"&sp=EgIQAQ%25253D%25253D",
		}

		for _, strategy := range searchStrategies {
			if len(videos) >= limit {
				break
			}

			altUrl := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query) + strategy
			altResp, err := httpClient.Get(altUrl)
			if err != nil {
				continue
			}

			altBody, err := io.ReadAll(altResp.Body)
			altResp.Body.Close()
			if err != nil {
				continue
			}

			altM := re.FindSubmatch(altBody)
			if len(altM) < 2 {
				continue
			}

			var altRoot map[string]interface{}
			if json.Unmarshal(altM[1], &altRoot) != nil {
				continue
			}

			altVideos := []types.Video{}
			var altWalk func(interface{})
			altWalk = func(node interface{}) {
				if m, ok := node.(map[string]interface{}); ok {
					if vr, ok := m["videoRenderer"]; ok {
						v := parseVideoRenderer(vr)
						if v.Title != "" && v.URL != "" {
							altVideos = append(altVideos, v)
						}
					}
					if cr, ok := m["compactVideoRenderer"]; ok {
						v := parseVideoRenderer(cr)
						if v.Title != "" && v.URL != "" {
							altVideos = append(altVideos, v)
						}
					}
					for _, v := range m {
						altWalk(v)
					}
				} else if arr, ok := node.([]interface{}); ok {
					for _, v := range arr {
						altWalk(v)
					}
				}
			}
			altWalk(altRoot)

			existingURLs := make(map[string]bool)
			for _, v := range videos {
				existingURLs[v.URL] = true
			}

			for _, v := range altVideos {
				if !existingURLs[v.URL] && len(videos) < limit {
					videos = append(videos, v)
					existingURLs[v.URL] = true
				}
			}
		}
	}

	if len(videos) > limit {
		videos = videos[:limit]
	}

	if progress != nil {
		progress(2, 2+len(videos))
	}

	total := len(videos)
	done := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(total)

	for i := range videos {
		go func(i int) {
			defer wg.Done()

			thumbPath := cacheThumbnailOptimized(videos[i].Thumbnail)
			if thumbPath == "" && videos[i].Thumbnail != "" {
				fallbackURLs := []string{
					strings.ReplaceAll(videos[i].Thumbnail, "default", "hqdefault"),
					strings.ReplaceAll(videos[i].Thumbnail, "default", "mqdefault"),
					strings.ReplaceAll(videos[i].Thumbnail, "default", "sddefault"),
					strings.ReplaceAll(videos[i].Thumbnail, "default", "maxresdefault"),
				}

				for _, fallbackURL := range fallbackURLs {
					if fallbackURL != videos[i].Thumbnail {
						thumbPath = cacheThumbnailOptimized(fallbackURL)
						if thumbPath != "" {
							break
						}
					}
				}
			}
			videos[i].ThumbnailPath = thumbPath

			mu.Lock()
			done++
			if progress != nil {
				progress(2+done, 2+total)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	CleanupHTTPConnections()

	return videos, nil
}

func cacheThumbnailOptimized(url string) string {
	if url == "" {
		return ""
	}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	cacheDir := "/tmp/gophertube-thumbs"
	os.MkdirAll(cacheDir, 0o755)
	thumbPath := filepath.Join(cacheDir, hash+".jpg")

	if _, err := os.Stat(thumbPath); err == nil {
		return thumbPath
	}

	for attempt := 0; attempt < 2; attempt++ {
		client := &http.Client{
			Timeout: time.Duration(3+attempt*2) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        0,
				MaxIdleConnsPerHost: 0,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression:  false,
				ForceAttemptHTTP2:   true,
				MaxConnsPerHost:     0,
				DisableKeepAlives:   false,
			},
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Referer", "https://www.youtube.com/")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			continue
		}

		data, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
		if err != nil || len(data) == 0 || isHTML(data) {
			continue
		}

		tempPath := thumbPath + ".tmp"
		if err := ioutil.WriteFile(tempPath, data, 0o644); err != nil {
			continue
		}

		if err := os.Rename(tempPath, thumbPath); err != nil {
			os.Remove(tempPath)
			continue
		}

		return thumbPath
	}

	return ""
}

// CleanupHTTPConnections closes idle connections to prevent memory leaks
func CleanupHTTPConnections() {
	if transport, ok := httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
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
	views := safeJQString(m, "viewCountText", "simpleText")
	thumb := ""
	if thumbs, ok := m["thumbnail"].(map[string]interface{}); ok {
		if arr, ok := thumbs["thumbnails"].([]interface{}); ok && len(arr) > 0 {
			for i := len(arr) - 1; i >= 0; i-- {
				if t, ok := arr[i].(map[string]interface{}); ok {
					if thumbURL, ok := t["url"].(string); ok && thumbURL != "" {
						if strings.Contains(thumbURL, "maxres") || strings.Contains(thumbURL, "hq") {
							thumb = thumbURL
							break
						} else if thumb == "" {
							thumb = thumbURL
						}
					}
				}
			}
		}
	}

	if thumb == "" {
		if thumbInfo, ok := m["thumbnailInfo"].(map[string]interface{}); ok {
			if thumbURL, ok := thumbInfo["url"].(string); ok && thumbURL != "" {
				thumb = thumbURL
			}
		}
	}
	return types.Video{
		Title:     title,
		URL:       url,
		Author:    channel,
		Duration:  duration,
		Views:     views,
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

func isHTML(data []byte) bool {
	return bytes.HasPrefix(data, []byte("<!DOCTYPE html")) || bytes.HasPrefix(data, []byte("<html"))
}
