package config

import (
	"fmt"
	"time"
)

// AggregatorConfig mirrors ports.AggregatorConfig for TOML/YAML unmarshalling.
type AggregatorConfig struct {
	Enabled  bool          `toml:"enabled"  yaml:"enabled"`
	Window   time.Duration `toml:"window"   yaml:"window"`
	MaxBatch int           `toml:"max_batch" yaml:"max_batch"`
}

// DefaultAggregatorConfig returns production-ready defaults.
func DefaultAggregatorConfig() AggregatorConfig {
	return AggregatorConfig{
		Enabled:  true,
		Window:   5 * time.Second,
		MaxBatch: 20,
	}
}

// Validate returns an error if the config contains invalid values.
func (c AggregatorConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return fmt.Errorf("aggregator: window must be positive, got %s", c.Window)
	}
	if c.MaxBatch < 1 {
		return fmt.Errorf("aggregator: max_batch must be at least 1, got %d", c.MaxBatch)
	}
	return nil
}
