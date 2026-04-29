package ports

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func makeSevPort(port uint16, sev Severity) Port {
	return Port{Port: port, Severity: sev}
}

func TestApplySeverityFilterKeepsAboveMin(t *testing.T) {
	cfg := config.SeverityFilterConfig{MinSeverity: SeverityWarning}
	ports := []Port{
		makeSevPort(80, SeverityInfo),
		makeSevPort(443, SeverityWarning),
		makeSevPort(4444, SeverityCritical),
	}
	got := ApplySeverityFilter(ports, cfg)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
	for _, p := range got {
		if p.Severity < SeverityWarning {
			t.Errorf("port %d has severity %v below minimum", p.Port, p.Severity)
		}
	}
}

func TestApplySeverityFilterAllPass(t *testing.T) {
	cfg := config.SeverityFilterConfig{MinSeverity: SeverityInfo}
	ports := []Port{
		makeSevPort(22, SeverityInfo),
		makeSevPort(8080, SeverityWarning),
	}
	got := ApplySeverityFilter(ports, cfg)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestApplySeverityFilterAllFiltered(t *testing.T) {
	cfg := config.SeverityFilterConfig{MinSeverity: SeverityCritical}
	ports := []Port{
		makeSevPort(80, SeverityInfo),
		makeSevPort(443, SeverityWarning),
	}
	got := ApplySeverityFilter(ports, cfg)
	if len(got) != 0 {
		t.Fatalf("expected 0 ports, got %d", len(got))
	}
}

func TestApplySeverityFilterEmpty(t *testing.T) {
	cfg := config.SeverityFilterConfig{MinSeverity: SeverityWarning}
	got := ApplySeverityFilter(nil, cfg)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestSeverityFilterFunc(t *testing.T) {
	cfg := config.SeverityFilterConfig{MinSeverity: SeverityWarning}
	f := SeverityFilter(cfg)
	if f(makeSevPort(80, SeverityInfo)) {
		t.Error("info port should be filtered")
	}
	if !f(makeSevPort(4444, SeverityCritical)) {
		t.Error("critical port should pass")
	}
}
