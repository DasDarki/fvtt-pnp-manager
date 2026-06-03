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
	body := map[string]any{"model": model, "prompt": prompt, "n": 1, "size": openAISize(model, size)}
	// gpt-image-* always return b64 and reject response_format; only DALL·E takes it.
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

// openAISize maps a requested WxH to a value the given model accepts.
// gpt-image-*: 1024x1024 | 1536x1024 | 1024x1536. dall-e-3: 1024x1024 | 1792x1024 | 1024x1792.
// dall-e-2: square only.
func openAISize(model, size string) string {
	wide := strings.HasPrefix(size, "1792x") || strings.HasPrefix(size, "1536x")
	tall := strings.HasSuffix(size, "x1792") || strings.HasSuffix(size, "x1536")
	switch {
	case strings.HasPrefix(model, "gpt-image"):
		if wide {
			return "1536x1024"
		}
		if tall {
			return "1024x1536"
		}
		return "1024x1024"
	case model == "dall-e-2":
		return "1024x1024"
	default: // dall-e-3 and unknown
		if wide {
			return "1792x1024"
		}
		if tall {
			return "1024x1792"
		}
		return "1024x1024"
	}
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
