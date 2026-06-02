package imagegen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func generateStability(ctx context.Context, cfg Config, prompt, size string) (Result, error) {
	segment := cfg.Model
	if segment == "" {
		segment = "core"
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("prompt", prompt)
	_ = w.WriteField("aspect_ratio", aspectRatio(size))
	_ = w.WriteField("output_format", "png")
	if segment == "sd3" {
		_ = w.WriteField("model", "sd3.5-large-turbo")
	}
	w.Close()

	url := "https://api.stability.ai/v2beta/stable-image/generate/" + segment
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Accept", "image/*")
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("stability %d: %s", resp.StatusCode, truncate(raw, 300))
	}
	return Result{Data: raw, Mime: "image/png", Ext: "png"}, nil
}
