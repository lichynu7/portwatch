package config

import (
	"testing"
	"time"
)

func TestDefaultWatcherConfig(t *testing.T) {
	cfg := DefaultWatcherConfig()

	if cfg.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %s", cfg.Interval)
	}
	if cfg.ProcFS != "/proc" {
		t.Errorf("expected /proc, got %q", cfg.ProcFS)
	}
	if cfg.EmitClosed {
		t.Error("expected EmitClosed to be false by default")
	}
	if cfg.BufferSize != 64 {
		t.Errorf("expected buffer size 64, got %d", cfg.BufferSize)
	}
}

func TestWatcherConfigValidateOK(t *testing.T) {
	cfg := DefaultWatcherConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestWatcherConfigValidateZeroInterval(t *testing.T) {
	cfg := DefaultWatcherConfig()
	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestWatcherConfigValidateNegativeInterval(t *testing.T) {
	cfg := DefaultWatcherConfig()
	cfg.Interval = -1 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestWatcherConfigValidateZeroBuffer(t *testing.T) {
	cfg := DefaultWatcherConfig()
	cfg.BufferSize = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero buffer size")
	}
}

func TestWatcherConfigValidateEmptyProcFS(t *testing.T) {
	cfg := DefaultWatcherConfig()
	cfg.ProcFS = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty proc_fs")
	}
}

func TestWatcherConfigEmitClosed(t *testing.T) {
	cfg := DefaultWatcherConfig()
	cfg.EmitClosed = true
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
