package ports

import (
	"testing"
	"time"
)

func fixedHour(hour int) func() time.Time {
	return func() time.Time {
		return time.Date(2024, 1, 1, hour, 0, 0, 0, time.UTC)
	}
}

func TestSuppressorNoQuietHours(t *testing.T) {
	cfg := DefaultSuppressorConfig()
	s := NewSuppressor(cfg)
	s.now = fixedHour(3)
	p := Port{Severity: SeverityInfo}
	if s.Suppress(p) {
		t.Error("expected no suppression when quiet hours disabled")
	}
}

func TestSuppressorWithinQuietHours(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: SeverityCritical}
	s := NewSuppressor(cfg)
	s.now = fixedHour(23)
	p := Port{Severity: SeverityInfo}
	if !s.Suppress(p) {
		t.Error("expected suppression inside quiet window")
	}
}

func TestSuppressorOutsideQuietHours(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: SeverityCritical}
	s := NewSuppressor(cfg)
	s.now = fixedHour(10)
	p := Port{Severity: SeverityInfo}
	if s.Suppress(p) {
		t.Error("expected no suppression outside quiet window")
	}
}

func TestSuppressorCriticalBypassesQuietHours(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 22, QuietHoursEnd: 6, MinSeverity: SeverityCritical}
	s := NewSuppressor(cfg)
	s.now = fixedHour(2)
	p := Port{Severity: SeverityCritical}
	if s.Suppress(p) {
		t.Error("critical alert should bypass quiet hours")
	}
}

func TestSuppressorDaytimeWindow(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 9, QuietHoursEnd: 17, MinSeverity: SeverityCritical}
	s := NewSuppressor(cfg)
	s.now = fixedHour(12)
	p := Port{Severity: SeverityWarning}
	if !s.Suppress(p) {
		t.Error("expected suppression inside daytime quiet window")
	}
}

func TestApplySuppressorNil(t *testing.T) {
	ports := []Port{{Port: 8080, Severity: SeverityInfo}}
	filter := ApplySuppressor(nil)
	out := filter(ports)
	if len(out) != 1 {
		t.Errorf("expected 1 port, got %d", len(out))
	}
}

func TestApplySuppressorFilters(t *testing.T) {
	cfg := SuppressorConfig{QuietHoursStart: 0, QuietHoursEnd: 23, MinSeverity: SeverityCritical}
	s := NewSuppressor(cfg)
	s.now = fixedHour(5)
	ports := []Port{
		{Port: 8080, Severity: SeverityInfo},
		{Port: 443, Severity: SeverityCritical},
	}
	filter := ApplySuppressor(s)
	out := filter(ports)
	if len(out) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out))
	}
	if out[0].Port != 443 {
		t.Errorf("expected critical port 443 to pass, got %d", out[0].Port)
	}
}
