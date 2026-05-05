package config

import (
	"testing"
)

func TestDefaultRateLimitFilterConfig(t *testing.T) {
	cfg := DefaultRateLimitFilterConfig()
	if !cfg.Enabled {
		t.Fatal("expected enabled by default")
	}
	if cfg.Window == "" {
		t.Fatal("expected non-empty window")
	}
	if cfg.MaxHits <= 0 {
		t.Fatalf("expected positive max_hits, got %d", cfg.MaxHits)
	}
}

func TestRateLimitFilterConfigValidateOK(t *testing.T) {
	cfg := DefaultRateLimitFilterConfig()
	if _, err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRateLimitFilterConfigValidateDisabled(t *testing.T) {
	cfg := RateLimitFilterConfig{Enabled: false}
	if _, err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should not error: %v", err)
	}
}

func TestRateLimitFilterConfigValidateEmptyWindow(t *testing.T) {
	cfg := RateLimitFilterConfig{Enabled: true, Window: "", MaxHits: 2}
	if _, err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty window")
	}
}

func TestRateLimitFilterConfigValidateInvalidWindow(t *testing.T) {
	cfg := RateLimitFilterConfig{Enabled: true, Window: "not-a-duration", MaxHits: 2}
	if _, err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestRateLimitFilterConfigValidateNegativeWindow(t *testing.T) {
	cfg := RateLimitFilterConfig{Enabled: true, Window: "-5s", MaxHits: 2}
	if _, err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestRateLimitFilterConfigValidateZeroMaxHits(t *testing.T) {
	cfg := RateLimitFilterConfig{Enabled: true, Window: "10s", MaxHits: 0}
	if _, err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero max_hits")
	}
}

func TestRateLimitFilterConfigFields(t *testing.T) {
	cfg := RateLimitFilterConfig{
		Enabled: true,
		Window:  "1m",
		MaxHits: 5,
	}
	d, err := cfg.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Seconds() != 60 {
		t.Fatalf("expected 60s window, got %v", d)
	}
}
