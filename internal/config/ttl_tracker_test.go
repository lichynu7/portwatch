package config

import (
	"testing"
	"time"
)

func TestDefaultTTLTrackerConfig(t *testing.T) {
	cfg := DefaultTTLTrackerConfig()
	if !cfg.Enabled {
		t.Fatal("expected enabled by default")
	}
	if cfg.TTL != 10*time.Minute {
		t.Fatalf("expected TTL 10m, got %v", cfg.TTL)
	}
	if cfg.EvictInterval != 2*time.Minute {
		t.Fatalf("expected EvictInterval 2m, got %v", cfg.EvictInterval)
	}
}

func TestTTLTrackerConfigValidateOK(t *testing.T) {
	cfg := DefaultTTLTrackerConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTTLTrackerConfigValidateDisabled(t *testing.T) {
	cfg := TTLTrackerConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestTTLTrackerConfigValidateZeroTTL(t *testing.T) {
	cfg := DefaultTTLTrackerConfig()
	cfg.TTL = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestTTLTrackerConfigValidateNegativeTTL(t *testing.T) {
	cfg := DefaultTTLTrackerConfig()
	cfg.TTL = -1 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestTTLTrackerConfigValidateZeroEvictInterval(t *testing.T) {
	cfg := DefaultTTLTrackerConfig()
	cfg.EvictInterval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero evict_interval")
	}
}

func TestTTLTrackerConfigValidateEvictIntervalExceedsTTL(t *testing.T) {
	cfg := TTLTrackerConfig{
		Enabled:       true,
		TTL:           1 * time.Minute,
		EvictInterval: 5 * time.Minute,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when evict_interval exceeds ttl")
	}
}
