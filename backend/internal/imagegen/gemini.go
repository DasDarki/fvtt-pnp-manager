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
	"time"
)

func generateGemini(ctx context.Context, cfg Config, prompt, size string) (Result, error) {
	model := cfg.Model
	if model == "" {
		model = "imagen-3.0-generate-002"
	}
	reqBody, _ := json.Marshal(map[string]any{
		"instances":  []map[string]any{{"prompt": prompt}},
		"parameters": map[string]any{"sampleCount": 1, "aspectRatio": aspectRatio(size)},
	})
	url := "https://generativelanguage.googleapis.com/v1beta/models/" + model + ":predict?key=" + cfg.APIKey
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("gemini %d: %s", resp.StatusCode, truncate(raw, 300))
	}

	var parsed struct {
		Predictions []struct {
			Bytes string `json:"bytesBase64Encoded"`
		} `json:"predictions"`
	}
	if json.Unmarshal(raw, &parsed) != nil || len(parsed.Predictions) == 0 || parsed.Predictions[0].Bytes == "" {
		return Result{}, errors.New("gemini: empty response")
	}
	img, derr := base64.StdEncoding.DecodeString(parsed.Predictions[0].Bytes)
	if derr != nil {
		return Result{}, derr
	}
	return Result{Data: img, Mime: "image/png", Ext: "png"}, nil
}
