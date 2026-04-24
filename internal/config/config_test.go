package config

import (
	"os"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.SnapshotDir != "/var/lib/portwatch" {
		t.Errorf("unexpected snapshot_dir: %s", cfg.SnapshotDir)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("unexpected log_level: %s", cfg.LogLevel)
	}
}

func TestLoadValid(t *testing.T) {
	path := writeTemp(t, "interval: 1m\nlog_level: debug\nwebhook:\n  url: http://hook.example.com\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != time.Minute {
		t.Errorf("expected 1m, got %v", cfg.Interval)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.LogLevel)
	}
	if cfg.Webhook.URL != "http://hook.example.com" {
		t.Errorf("unexpected webhook url: %s", cfg.Webhook.URL)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadInvalidInterval(t *testing.T) {
	path := writeTemp(t, "interval: -5s\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestLoadWebhookTimeout(t *testing.T) {
	path := writeTemp(t, "interval: 10s\nwebhook:\n  url: http://x.example.com\n  timeout: 10s\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Webhook.Timeout != 10*time.Second {
		t.Errorf("expected 10s webhook timeout, got %v", cfg.Webhook.Timeout)
	}
}
