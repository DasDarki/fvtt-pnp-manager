package imagegen

import (
	"context"
	"errors"
	"strings"
)

type Config struct {
	Provider string
	Model    string
	APIKey   string
}

type Result struct {
	Data []byte
	Mime string
	Ext  string
	Mock bool
}

// Providers are the BYOK image backends the user can choose.
var Providers = []string{"openai", "stability", "replicate", "gemini"}

func Generate(ctx context.Context, cfg Config, prompt, size string) (Result, error) {
	if cfg.APIKey == "" || cfg.Provider == "" {
		return Result{Data: []byte(mockSVG(prompt, size)), Mime: "image/svg+xml", Ext: "svg", Mock: true}, nil
	}
	switch cfg.Provider {
	case "openai":
		return generateOpenAI(ctx, cfg, prompt, size)
	case "stability":
		return generateStability(ctx, cfg, prompt, size)
	case "replicate":
		return generateReplicate(ctx, cfg, prompt, size)
	case "gemini":
		return generateGemini(ctx, cfg, prompt, size)
	default:
		return Result{}, errors.New("unknown image provider: " + cfg.Provider)
	}
}

func aspectRatio(size string) string {
	switch {
	case strings.HasPrefix(size, "1792x") || strings.HasPrefix(size, "1536x"):
		return "16:9"
	case strings.HasSuffix(size, "x1792") || strings.HasSuffix(size, "x1536"):
		return "9:16"
	default:
		return "1:1"
	}
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) > n {
		return s[:n]
	}
	return s
}
