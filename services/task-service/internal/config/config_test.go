package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	withTempDir(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTP.Addr != ":8082" {
		t.Fatalf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, ":8082")
	}
	if cfg.DB.URL != defaultDBURL {
		t.Fatalf("DB.URL = %q, want %q", cfg.DB.URL, defaultDBURL)
	}
}

func TestLoadFileOverridesDefaults(t *testing.T) {
	dir := withTempDir(t)

	config := `
http:
  addr: ":9090"
db:
  url: "postgres://example.com:5432/app?sslmode=disable"
`
	if err := os.WriteFile(dir+"/config.yaml", []byte(config), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTP.Addr != ":9090" {
		t.Fatalf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, ":9090")
	}
	if cfg.DB.URL != "postgres://example.com:5432/app?sslmode=disable" {
		t.Fatalf("DB.URL = %q, want file value", cfg.DB.URL)
	}
}

func TestLoadEnvOverridesFile(t *testing.T) {
	dir := withTempDir(t)

	config := `
http:
  addr: ":9090"
db:
  url: "postgres://example.com:5432/app?sslmode=disable"
`
	if err := os.WriteFile(dir+"/config.yaml", []byte(config), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	t.Setenv("TASKBOARD_HTTP_ADDR", ":9191")
	t.Setenv("TASKBOARD_DB_URL", "postgres://env.example.com:5432/env?sslmode=disable")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTP.Addr != ":9191" {
		t.Fatalf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, ":9191")
	}
	if cfg.DB.URL != "postgres://env.example.com:5432/env?sslmode=disable" {
		t.Fatalf("DB.URL = %q, want env value", cfg.DB.URL)
	}
}

func TestLoadValidation(t *testing.T) {
	withTempDir(t)

	t.Setenv("TASKBOARD_DB_URL", "   ")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "db.url") {
		t.Fatalf("Load() error = %q, want db.url mention", err.Error())
	}
}

func withTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	return dir
}
