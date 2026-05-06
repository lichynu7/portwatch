package config

import "testing"

func TestDefaultChangelogConfig(t *testing.T) {
	cfg := DefaultChangelogConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if cfg.MaxEvents != 256 {
		t.Errorf("expected MaxEvents 256, got %d", cfg.MaxEvents)
	}
}

func TestChangelogConfigValidateOK(t *testing.T) {
	cfg := DefaultChangelogConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestChangelogConfigValidateDisabled(t *testing.T) {
	cfg := ChangelogConfig{Enabled: false, MaxEvents: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got: %v", err)
	}
}

func TestChangelogConfigValidateZeroMaxEvents(t *testing.T) {
	cfg := ChangelogConfig{Enabled: true, MaxEvents: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for MaxEvents=0")
	}
}

func TestChangelogConfigValidateNegativeMaxEvents(t *testing.T) {
	cfg := ChangelogConfig{Enabled: true, MaxEvents: -5}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative MaxEvents")
	}
}

func TestChangelogConfigValidateExceedsMax(t *testing.T) {
	cfg := ChangelogConfig{Enabled: true, MaxEvents: 99999}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for MaxEvents exceeding limit")
	}
}

func TestChangelogConfigFields(t *testing.T) {
	cfg := ChangelogConfig{Enabled: true, MaxEvents: 512}
	if cfg.MaxEvents != 512 {
		t.Errorf("expected 512, got %d", cfg.MaxEvents)
	}
}
