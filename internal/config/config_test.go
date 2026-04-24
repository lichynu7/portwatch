package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %s", cfg.Interval)
	}
	if !cfg.Alerts.Stdout {
		t.Error("expected stdout alerts enabled by default")
	}
}

func TestLoadValid(t *testing.T) {
	yaml := `
interval: 10s
snapshot_dir: /tmp/pw
allowed_ports: [22, 80, 443]
alerts:
  stdout: true
  log_file: /tmp/pw.log
`
	p := writeTemp(t, yaml)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("interval: got %s, want 10s", cfg.Interval)
	}
	if len(cfg.AllowedPorts) != 3 {
		t.Errorf("allowed_ports: got %d entries, want 3", len(cfg.AllowedPorts))
	}
	if cfg.Alerts.LogFile != "/tmp/pw.log" {
		t.Errorf("log_file: got %q", cfg.Alerts.LogFile)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadInvalidInterval(t *testing.T) {
	yaml := `interval: 500ms\nsnapshot_dir: /tmp/pw\n`
	p := writeTemp(t, yaml)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for sub-second interval")
	}
}

func TestLoadUnknownField(t *testing.T) {
	yaml := `unknown_key: oops\ninterval: 5s\nsnapshot_dir: /tmp/pw\n`
	p := writeTemp(t, yaml)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}
