package config

import (
	"fmt"
	"time"
)

// TTLTrackerConfig holds configuration for the TTLTracker.
type TTLTrackerConfig struct {
	// Enabled controls whether TTL tracking is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// TTL is the duration after which an unseen port entry is evicted.
	TTL time.Duration `toml:"ttl" yaml:"ttl"`

	// EvictInterval is how often the eviction sweep runs.
	EvictInterval time.Duration `toml:"evict_interval" yaml:"evict_interval"`
}

// DefaultTTLTrackerConfig returns a TTLTrackerConfig with sensible defaults.
func DefaultTTLTrackerConfig() TTLTrackerConfig {
	return TTLTrackerConfig{
		Enabled:       true,
		TTL:           10 * time.Minute,
		EvictInterval: 2 * time.Minute,
	}
}

// Validate checks that the TTLTrackerConfig fields are valid.
func (c TTLTrackerConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TTL <= 0 {
		return fmt.Errorf("ttl_tracker: ttl must be positive, got %v", c.TTL)
	}
	if c.EvictInterval <= 0 {
		return fmt.Errorf("ttl_tracker: evict_interval must be positive, got %v", c.EvictInterval)
	}
	if c.EvictInterval > c.TTL {
		return fmt.Errorf("ttl_tracker: evict_interval (%v) should not exceed ttl (%v)", c.EvictInterval, c.TTL)
	}
	return nil
}
