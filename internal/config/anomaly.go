package config

import (
	"fmt"
	"time"
)

// AnomalyConfig holds tuning parameters for the anomaly detector.
type AnomalyConfig struct {
	// Enabled controls whether anomaly detection is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// Window is the rolling duration within which occurrences are counted.
	Window time.Duration `toml:"window" yaml:"window"`

	// MinOccurrences is the minimum number of times a port must be observed
	// within Window before it is flagged as an anomaly.
	MinOccurrences int `toml:"min_occurrences" yaml:"min_occurrences"`
}

// DefaultAnomalyConfig returns a sensible out-of-the-box AnomalyConfig.
func DefaultAnomalyConfig() AnomalyConfig {
	return AnomalyConfig{
		Enabled:        true,
		Window:         2 * time.Minute,
		MinOccurrences: 3,
	}
}

// Validate returns an error if the configuration is invalid.
func (c AnomalyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return fmt.Errorf("anomaly: window must be positive, got %s", c.Window)
	}
	if c.MinOccurrences < 1 {
		return fmt.Errorf("anomaly: min_occurrences must be >= 1, got %d", c.MinOccurrences)
	}
	return nil
}
