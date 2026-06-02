package imagegen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type replicatePrediction struct {
	Status string          `json:"status"`
	Output json.RawMessage `json:"output"`
	Error  string          `json:"error"`
	URLs   struct {
		Get string `json:"get"`
	} `json:"urls"`
}

func generateReplicate(ctx context.Context, cfg Config, prompt, size string) (Result, error) {
	model := cfg.Model
	if model == "" {
		model = "black-forest-labs/flux-schnell"
	}
	reqBody, _ := json.Marshal(map[string]any{
		"input": map[string]any{"prompt": prompt, "aspect_ratio": aspectRatio(size), "output_format": "png"},
	})
	url := "https://api.replicate.com/v1/models/" + model + "/predictions"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "wait")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("replicate %d: %s", resp.StatusCode, truncate(raw, 300))
	}

	var pred replicatePrediction
	if json.Unmarshal(raw, &pred) != nil {
		return Result{}, errors.New("replicate: bad response")
	}

	for i := 0; i < 30 && pred.Status != "succeeded" && pred.Status != "failed" && pred.URLs.Get != ""; i++ {
		time.Sleep(1500 * time.Millisecond)
		greq, _ := http.NewRequestWithContext(ctx, http.MethodGet, pred.URLs.Get, nil)
		greq.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		gresp, gerr := client.Do(greq)
		if gerr != nil {
			break
		}
		body, _ := io.ReadAll(gresp.Body)
		gresp.Body.Close()
		_ = json.Unmarshal(body, &pred)
	}

	if pred.Status == "failed" {
		return Result{}, fmt.Errorf("replicate failed: %s", pred.Error)
	}
	out := firstURL(pred.Output)
	if out == "" {
		return Result{}, errors.New("replicate: no output")
	}
	return downloadImage(ctx, out)
}

func firstURL(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil && s != "" {
		return s
	}
	var arr []string
	if json.Unmarshal(raw, &arr) == nil && len(arr) > 0 {
		return arr[0]
	}
	return ""
}
