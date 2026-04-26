package config

import (
	"fmt"
	"time"
)

// DedupConfig holds TOML-serialisable deduplication settings.
type DedupConfig struct {
	// WindowSize is the duration a port alert is suppressed after first firing.
	// Accepts Go duration strings, e.g. "5m", "30s".
	WindowSize string `toml:"window_size"`
}

// DefaultDedupConfig returns production-safe defaults.
func DefaultDedupConfig() DedupConfig {
	return DedupConfig{
		WindowSize: "5m",
	}
}

// WindowDuration parses WindowSize into a time.Duration.
func (c DedupConfig) WindowDuration() (time.Duration, error) {
	if c.WindowSize == "" {
		return 0, fmt.Errorf("dedup window_size must not be empty")
	}
	d, err := time.ParseDuration(c.WindowSize)
	if err != nil {
		return 0, fmt.Errorf("dedup window_size %q is not a valid duration: %w", c.WindowSize, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("dedup window_size must be positive, got %q", c.WindowSize)
	}
	return d, nil
}

// Validate returns an error if the configuration is invalid.
func (c DedupConfig) Validate() error {
	_, err := c.WindowDuration()
	return err
}
