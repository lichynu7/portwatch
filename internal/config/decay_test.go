package config

import (
	"testing"
	"time"
)

func TestDefaultDecayConfig(t *testing.T) {
	cfg := DefaultDecayConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.HalfLife != 10*time.Minute {
		t.Errorf("expected 10m half-life, got %s", cfg.HalfLife)
	}
	if cfg.MinScore != 1.0 {
		t.Errorf("expected MinScore 1.0, got %f", cfg.MinScore)
	}
}

func TestDecayConfigValidateOK(t *testing.T) {
	cfg := DefaultDecayConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDecayConfigValidateDisabled(t *testing.T) {
	cfg := DecayConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestDecayConfigValidateZeroHalfLife(t *testing.T) {
	cfg := DefaultDecayConfig()
	cfg.HalfLife = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero half_life")
	}
}

func TestDecayConfigValidateNegativeHalfLife(t *testing.T) {
	cfg := DefaultDecayConfig()
	cfg.HalfLife = -1 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative half_life")
	}
}

func TestDecayConfigValidateNegativeMinScore(t *testing.T) {
	cfg := DefaultDecayConfig()
	cfg.MinScore = -0.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative min_score")
	}
}

func TestDecayConfigValidateZeroMinScoreOK(t *testing.T) {
	cfg := DefaultDecayConfig()
	cfg.MinScore = 0
	if err := cfg.Validate(); err != nil {
		t.Fatalf("zero min_score should be valid, got: %v", err)
	}
}
