package imagegen

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func generateOpenAI(ctx context.Context, cfg Config, prompt, size string) (Result, error) {
	model := cfg.Model
	if model == "" {
		model = "dall-e-3"
	}
	body := map[string]any{"model": model, "prompt": prompt, "n": 1, "size": size}
	if strings.HasPrefix(model, "dall-e") {
		body["response_format"] = "b64_json"
	}
	buf, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/images/generations", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("openai %d: %s", resp.StatusCode, truncate(raw, 300))
	}

	var parsed struct {
		Data []struct {
			B64 string `json:"b64_json"`
			URL string `json:"url"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &parsed) != nil || len(parsed.Data) == 0 {
		return Result{}, errors.New("openai: empty response")
	}
	if parsed.Data[0].B64 != "" {
		img, derr := base64.StdEncoding.DecodeString(parsed.Data[0].B64)
		if derr != nil {
			return Result{}, derr
		}
		return Result{Data: img, Mime: "image/png", Ext: "png"}, nil
	}
	if parsed.Data[0].URL != "" {
		return downloadImage(ctx, parsed.Data[0].URL)
	}
	return Result{}, errors.New("openai: no image data")
}

func downloadImage(ctx context.Context, url string) (Result, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("download %d", resp.StatusCode)
	}
	img, _ := io.ReadAll(resp.Body)
	mime := resp.Header.Get("Content-Type")
	ext := "png"
	if strings.Contains(mime, "webp") {
		ext = "webp"
	} else if strings.Contains(mime, "jpeg") {
		ext = "jpg"
	}
	if mime == "" {
		mime = "image/png"
	}
	return Result{Data: img, Mime: mime, Ext: ext}, nil
}
