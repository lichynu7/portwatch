package config

import "testing"

func TestDefaultSuppressorConfig(t *testing.T) {
	cfg := DefaultSuppressorConfig()
	if cfg.QuietHoursStart != -1 || cfg.QuietHoursEnd != -1 {
		t.Error("default suppressor should have quiet hours disabled")
	}
	if cfg.MinSeverity != "critical" {
		t.Errorf("expected min_severity=critical, got %s", cfg.MinSeverity)
	}
}

func TestSuppressorConfigValidateDisabled(t *testing.T) {
	cfg := DefaultSuppressorConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for disabled config: %v", err)
	}
}

func TestSuppressorConfigValidateOK(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: "critical"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSuppressorConfigValidateInvalidStart(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 25, QuietHoursEnd: 6, MinSeverity: "critical"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for out-of-range start hour")
	}
}

func TestSuppressorConfigValidateInvalidEnd(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: -2, MinSeverity: "critical"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for out-of-range end hour")
	}
}

func TestSuppressorConfigValidateInvalidSeverity(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: "unknown"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unknown severity")
	}
}

func TestSuppressorConfigToPortsConfig(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: "warning"}
	pc := cfg.ToPortsSuppressorConfig()
	if pc.QuietHoursStart != 22 {
		t.Errorf("expected start=22, got %d", pc.QuietHoursStart)
	}
	if pc.QuietHoursEnd != 6 {
		t.Errorf("expected end=6, got %d", pc.QuietHoursEnd)
	}
}
