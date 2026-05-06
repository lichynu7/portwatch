package config

import (
	"fmt"
	"time"
)

// DecayConfig holds configuration for the score-decay feature.
type DecayConfig struct {
	Enabled  bool          `toml:"enabled" yaml:"enabled"`
	HalfLife time.Duration `toml:"half_life" yaml:"half_life"`
	MinScore float64       `toml:"min_score" yaml:"min_score"`
}

// DefaultDecayConfig returns production-ready defaults.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		Enabled:  true,
		HalfLife: 10 * time.Minute,
		MinScore: 1.0,
	}
}

// Validate returns an error if the configuration is invalid.
func (c DecayConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.HalfLife <= 0 {
		return fmt.Errorf("decay: half_life must be positive, got %s", c.HalfLife)
	}
	if c.MinScore < 0 {
		return fmt.Errorf("decay: min_score must be >= 0, got %f", c.MinScore)
	}
	return nil
}
