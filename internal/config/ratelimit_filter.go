package config

import (
	"fmt"
	"time"
)

// RateLimitFilterConfig holds the configuration for the alert-rate-limit
// pipeline stage. It is distinct from the per-port RateLimiter and operates
// at the pipeline level, capping burst throughput for repeated port alerts.
type RateLimitFilterConfig struct {
	// Enabled controls whether the stage is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`
	// Window is the rolling time window expressed as a duration string.
	Window string `toml:"window" yaml:"window"`
	// MaxHits is the maximum number of alerts allowed per port key per window.
	MaxHits int `toml:"max_hits" yaml:"max_hits"`
}

// DefaultRateLimitFilterConfig returns production-ready defaults.
func DefaultRateLimitFilterConfig() RateLimitFilterConfig {
	return RateLimitFilterConfig{
		Enabled: true,
		Window:  "30s",
		MaxHits: 3,
	}
}

// Validate checks the config for logical errors and returns the parsed window
// duration on success.
func (c RateLimitFilterConfig) Validate() (time.Duration, error) {
	if !c.Enabled {
		return 0, nil
	}
	if c.Window == "" {
		return 0, fmt.Errorf("ratelimit_filter: window must not be empty")
	}
	d, err := time.ParseDuration(c.Window)
	if err != nil {
		return 0, fmt.Errorf("ratelimit_filter: invalid window %q: %w", c.Window, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("ratelimit_filter: window must be positive")
	}
	if c.MaxHits <= 0 {
		return 0, fmt.Errorf("ratelimit_filter: max_hits must be > 0")
	}
	return d, nil
}
