package config

import (
	"os"

	"github.com/subosito/gotenv"
)

// Config holds runtime configuration for the task service.
type Config struct {
	Addr        string
	DatabaseURL string
}

// FromEnv builds a Config using environment variables with sensible defaults.
// It attempts to load .env files from common locations so that `go run` works
// whether you execute from the repo root or the service directory.
func FromEnv() Config {
	loadEnvFiles()

	return Config{
		Addr:        envOrDefault("TASK_HTTP_ADDR", ":8082"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func loadEnvFiles() {
	paths := []string{
		".env",
		"services/task-service/.env",
		"../.env",
		"../../.env",
	}
	for _, p := range paths {
		_ = gotenv.Load(p)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
