package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	CORSOrigins   string
	OpenAIKey     string
	OpenAIModel   string
	UploadDir     string
	PublicBaseURL string
	EncryptionKey string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Port:          env("PORT", "8080"),
		DatabaseURL:   env("DATABASE_URL", "postgres://aetherwright:aetherwright@localhost:5432/aetherwright?sslmode=disable"),
		JWTSecret:     env("JWT_SECRET", "dev-secret-change-me"),
		AccessTTL:     envDuration("ACCESS_TTL", 15*time.Minute),
		RefreshTTL:    envDuration("REFRESH_TTL", 30*24*time.Hour),
		CORSOrigins:   env("CORS_ORIGINS", "http://localhost:3000"),
		OpenAIKey:     env("OPENAI_API_KEY", ""),
		OpenAIModel:   env("OPENAI_IMAGE_MODEL", "dall-e-3"),
		UploadDir:     env("UPLOAD_DIR", "./uploads"),
		PublicBaseURL: env("PUBLIC_BASE_URL", "http://localhost:8080"),
		EncryptionKey: env("ENCRYPTION_KEY", env("JWT_SECRET", "dev-secret-change-me")),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
