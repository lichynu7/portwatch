package ports

import (
	"sync"
	"time"
)

// Score represents a composite risk score for a port observation.
type Score struct {
	Port     uint16
	Protocol string
	Value    float64 // 0.0 (benign) – 1.0 (critical)
	Reasons  []string
	ScoredAt time.Time
}

// ScorerConfig controls how individual signals are weighted.
type ScorerConfig struct {
	SeverityWeight   float64 // contribution from severity level
	AnomalyWeight    float64 // contribution when anomaly flag is set
	EscalationWeight float64 // contribution when escalation flag is set
	GeoWeight        float64 // contribution for public/external IPs
	MaxScore         float64 // cap (default 1.0)
}

// DefaultScorerConfig returns sensible defaults.
func DefaultScorerConfig() ScorerConfig {
	return ScorerConfig{
		SeverityWeight:   0.40,
		AnomalyWeight:    0.25,
		EscalationWeight: 0.20,
		GeoWeight:        0.15,
		MaxScore:         1.0,
	}
}

// Scorer computes a composite risk score for each port.
type Scorer struct {
	cfg ScorerConfig
	mu  sync.Mutex
	last map[string]Score
}

// NewScorer creates a Scorer with the provided config.
func NewScorer(cfg ScorerConfig) *Scorer {
	if cfg.MaxScore <= 0 {
		cfg.MaxScore = 1.0
	}
	return &Scorer{
		cfg:  cfg,
		last: make(map[string]Score),
	}
}

// Score computes and caches the risk score for p.
func (s *Scorer) Score(p Port) Score {
	var total float64
	var reasons []string

	// Severity contribution.
	sevNorm := float64(p.Severity) / float64(SeverityCritical)
	if contribution := sevNorm * s.cfg.SeverityWeight; contribution > 0 {
		total += contribution
		reasons = append(reasons, "severity")
	}

	// Anomaly flag.
	if p.Anomaly {
		total += s.cfg.AnomalyWeight
		reasons = append(reasons, "anomaly")
	}

	// Escalation flag.
	if p.Escalated {
		total += s.cfg.EscalationWeight
		reasons = append(reasons, "escalation")
	}

	// Geo: public IP.
	if p.GeoScope == "public" {
		total += s.cfg.GeoWeight
		reasons = append(reasons, "public-ip")
	}

	if total > s.cfg.MaxScore {
		total = s.cfg.MaxScore
	}

	sc := Score{
		Port:     p.Port,
		Protocol: p.Protocol,
		Value:    total,
		Reasons:  reasons,
		ScoredAt: time.Now(),
	}

	s.mu.Lock()
	s.last[portKey(p)] = sc
	s.mu.Unlock()
	return sc
}

// Last returns the most recently computed score for the given port key.
func (s *Scorer) Last(key string) (Score, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sc, ok := s.last[key]
	return sc, ok
}
