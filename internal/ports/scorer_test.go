package ports

import (
	"testing"
)

func basePort(port uint16, sev Severity) Port {
	return Port{
		Port:     port,
		Protocol: "tcp",
		Severity: sev,
	}
}

func TestScorerBenignPort(t *testing.T) {
	s := NewScorer(DefaultScorerConfig())
	p := basePort(80, SeverityInfo)
	sc := s.Score(p)
	if sc.Value != 0 {
		t.Fatalf("expected 0 score for benign port, got %v", sc.Value)
	}
	if len(sc.Reasons) != 0 {
		t.Fatalf("expected no reasons, got %v", sc.Reasons)
	}
}

func TestScorerSeverityContribution(t *testing.T) {
	s := NewScorer(DefaultScorerConfig())
	p := basePort(4444, SeverityCritical)
	sc := s.Score(p)
	// severity contribution = 1.0 * 0.40 = 0.40
	if sc.Value < 0.39 || sc.Value > 0.41 {
		t.Fatalf("expected ~0.40, got %v", sc.Value)
	}
	if !containsReason(sc.Reasons, "severity") {
		t.Fatal("expected 'severity' in reasons")
	}
}

func TestScorerAnomalyFlag(t *testing.T) {
	s := NewScorer(DefaultScorerConfig())
	p := basePort(9999, SeverityInfo)
	p.Anomaly = true
	sc := s.Score(p)
	if sc.Value < 0.24 || sc.Value > 0.26 {
		t.Fatalf("expected ~0.25, got %v", sc.Value)
	}
	if !containsReason(sc.Reasons, "anomaly") {
		t.Fatal("expected 'anomaly' in reasons")
	}
}

func TestScorerCapsAtMaxScore(t *testing.T) {
	cfg := DefaultScorerConfig()
	cfg.MaxScore = 0.5
	s := NewScorer(cfg)
	p := basePort(1234, SeverityCritical)
	p.Anomaly = true
	p.Escalated = true
	p.GeoScope = "public"
	sc := s.Score(p)
	if sc.Value > 0.5 {
		t.Fatalf("score %v exceeds max_score 0.5", sc.Value)
	}
}

func TestScorerPublicGeo(t *testing.T) {
	s := NewScorer(DefaultScorerConfig())
	p := basePort(8080, SeverityInfo)
	p.GeoScope = "public"
	sc := s.Score(p)
	if !containsReason(sc.Reasons, "public-ip") {
		t.Fatal("expected 'public-ip' in reasons")
	}
}

func TestScorerLastCached(t *testing.T) {
	s := NewScorer(DefaultScorerConfig())
	p := basePort(22, SeverityWarning)
	sc := s.Score(p)
	cached, ok := s.Last(portKey(p))
	if !ok {
		t.Fatal("expected cached score")
	}
	if cached.Value != sc.Value {
		t.Fatalf("cached value %v != scored value %v", cached.Value, sc.Value)
	}
}

func containsReason(reasons []string, target string) bool {
	for _, r := range reasons {
		if r == target {
			return true
		}
	}
	return false
}
