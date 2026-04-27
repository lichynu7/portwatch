package config

import "testing"

func TestDefaultSeverityFilterConfig(t *testing.T) {
	cfg := DefaultSeverityFilterConfig()
	if cfg.MinLevel != "info" {
		t.Errorf("expected default min_level \"info\", got %q", cfg.MinLevel)
	}
}

func TestSeverityFilterConfigValidateOK(t *testing.T) {
	for _, level := range []string{"info", "warning", "critical"} {
		cfg := SeverityFilterConfig{MinLevel: level}
		if err := cfg.Validate(); err != nil {
			t.Errorf("level %q: unexpected error: %v", level, err)
		}
	}
}

func TestSeverityFilterConfigValidateInvalid(t *testing.T) {
	cfg := SeverityFilterConfig{MinLevel: "debug"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unknown level \"debug\", got nil")
	}
}

func TestSeverityFilterConfigValidateEmpty(t *testing.T) {
	cfg := SeverityFilterConfig{MinLevel: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty min_level, got nil")
	}
}
