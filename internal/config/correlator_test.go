package config

import (
	"testing"
)

func TestDefaultCorrelatorConfig(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	if !cfg.Enabled {
		t.Error("expected correlator to be enabled by default")
	}
	if cfg.WindowDuration != "30s" {
		t.Errorf("expected window_duration 30s, got %s", cfg.WindowDuration)
	}
	if cfg.MinOccurrences != 2 {
		t.Errorf("expected min_occurrences 2, got %d", cfg.MinOccurrences)
	}
}

func TestCorrelatorConfigValidateOK(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestCorrelatorConfigValidateDisabled(t *testing.T) {
	cfg := CorrelatorConfig{Enabled: false, WindowDuration: "", MinOccurrences: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got: %v", err)
	}
}

func TestCorrelatorConfigValidateInvalidDuration(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	cfg.WindowDuration = "not-a-duration"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid window_duration")
	}
}

func TestCorrelatorConfigValidateZeroDuration(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	cfg.WindowDuration = "0s"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero window_duration")
	}
}

func TestCorrelatorConfigValidateZeroMinOccurrences(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	cfg.MinOccurrences = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for min_occurrences < 1")
	}
}

func TestCorrelatorConfigWindowDurationParsed(t *testing.T) {
	cfg := DefaultCorrelatorConfig()
	d, err := cfg.WindowDurationParsed()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Seconds() != 30 {
		t.Errorf("expected 30s, got %v", d)
	}
}
