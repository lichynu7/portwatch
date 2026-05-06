package config

import (
	"fmt"
	"time"
)

// SamplerConfig mirrors ports.SamplerConfig for TOML/JSON unmarshalling.
type SamplerConfig struct {
	Enabled       bool          `toml:"enabled"        json:"enabled"`
	ReservoirSize int           `toml:"reservoir_size" json:"reservoir_size"`
	Window        time.Duration `toml:"window"         json:"window"`
}

// DefaultSamplerConfig returns production defaults.
func DefaultSamplerConfig() SamplerConfig {
	return SamplerConfig{
		Enabled:       false,
		ReservoirSize: 100,
		Window:        time.Minute,
	}
}

// Validate returns an error if the configuration is inconsistent.
func (c SamplerConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.ReservoirSize <= 0 {
		return fmt.Errorf("sampler: reservoir_size must be > 0, got %d", c.ReservoirSize)
	}
	if c.Window <= 0 {
		return fmt.Errorf("sampler: window must be > 0, got %s", c.Window)
	}
	return nil
}
