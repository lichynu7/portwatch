package config

import (
	"testing"
	"time"
)

func TestDefaultDedupConfig(t *testing.T) {
	cfg := DefaultDedupConfig()
	if cfg.WindowSize == "" {
		t.Fatal("expected non-empty default WindowSize")
	}
	d, err := cfg.WindowDuration()
	if err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
	if d <= 0 {
		t.Fatalf("expected positive duration, got %v", d)
	}
}

func TestDedupConfigValidateOK(t *testing.T) {
	cfg := DedupConfig{WindowSize: "10m"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDedupConfigValidateEmpty(t *testing.T) {
	cfg := DedupConfig{WindowSize: ""}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty WindowSize")
	}
}

func TestDedupConfigValidateInvalid(t *testing.T) {
	cfg := DedupConfig{WindowSize: "not-a-duration"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid duration string")
	}
}

func TestDedupConfigValidateNegative(t *testing.T) {
	cfg := DedupConfig{WindowSize: "-1m"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestDedupConfigWindowDuration(t *testing.T) {
	cfg := DedupConfig{WindowSize: "2m30s"}
	d, err := cfg.WindowDuration()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 2*time.Minute + 30*time.Second
	if d != expected {
		t.Fatalf("expected %v, got %v", expected, d)
	}
}
