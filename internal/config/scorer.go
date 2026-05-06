package config

import "fmt"

// ScorerConfig mirrors ports.ScorerConfig for TOML/JSON unmarshalling.
type ScorerConfig struct {
	SeverityWeight   float64 `toml:"severity_weight"   json:"severity_weight"`
	AnomalyWeight    float64 `toml:"anomaly_weight"    json:"anomaly_weight"`
	EscalationWeight float64 `toml:"escalation_weight" json:"escalation_weight"`
	GeoWeight        float64 `toml:"geo_weight"        json:"geo_weight"`
	MaxScore         float64 `toml:"max_score"         json:"max_score"`
}

// DefaultScorerConfig returns sensible defaults aligned with the port scorer.
func DefaultScorerConfig() ScorerConfig {
	return ScorerConfig{
		SeverityWeight:   0.40,
		AnomalyWeight:    0.25,
		EscalationWeight: 0.20,
		GeoWeight:        0.15,
		MaxScore:         1.0,
	}
}

// Validate returns an error if any weight is negative or the weights exceed MaxScore.
func (c ScorerConfig) Validate() error {
	if c.MaxScore <= 0 {
		return fmt.Errorf("scorer: max_score must be > 0, got %v", c.MaxScore)
	}
	weights := []struct {
		name string
		v    float64
	}{
		{"severity_weight", c.SeverityWeight},
		{"anomaly_weight", c.AnomalyWeight},
		{"escalation_weight", c.EscalationWeight},
		{"geo_weight", c.GeoWeight},
	}
	var total float64
	for _, w := range weights {
		if w.v < 0 {
			return fmt.Errorf("scorer: %s must be >= 0, got %v", w.name, w.v)
		}
		total += w.v
	}
	if total > c.MaxScore {
		return fmt.Errorf("scorer: sum of weights (%.2f) exceeds max_score (%.2f)", total, c.MaxScore)
	}
	return nil
}
