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
)

func SearchYouTube(query string, limit int) ([]types.Video, error) {
	url := "https://www.youtube.com/results?search_query=" + urlQueryEscape(query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
		if len(videos) >= limit {
			return
		}
		if m, ok := node.(map[string]interface{}); ok {
			if vr, ok := m["videoRenderer"]; ok {
				v := parseVideoRenderer(vr)
				if v.Title != "" && v.URL != "" {
					v.ThumbnailPath = cacheThumbnail(v.Thumbnail)
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
	if len(videos) > limit {
		videos = videos[:limit]
	}
	return videos, nil
}

func cacheThumbnail(url string) string {
	if url == "" {
		return ""
	}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	cacheDir := "/tmp/gophertube-thumbs"
	os.MkdirAll(cacheDir, 0o755)
	thumbPath := filepath.Join(cacheDir, hash+".jpg")
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GopherTube/1.0; +https://github.com/KrishnaSSH/GopherTube)")
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err == nil && len(data) > 0 && !isHTML(data) {
				ioutil.WriteFile(thumbPath, data, 0o644)
			}
		}
	}
	return thumbPath
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
