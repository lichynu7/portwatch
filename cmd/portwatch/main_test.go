package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestLoadConfigNoPath(t *testing.T) {
	cfg, err := loadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := config.Default()
	if cfg.Interval != def.Interval {
		t.Errorf("expected default interval %s, got %s", def.Interval, cfg.Interval)
	}
}

func TestLoadConfigValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.toml")
	content := `interval = "30s"
snapshot_path = "/tmp/snap.json"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval.String() != "30s" {
		t.Errorf("expected interval 30s, got %s", cfg.Interval)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/portwatch.toml")
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}
