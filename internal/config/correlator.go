package config

import (
	"fmt"
	"time"
)

// CorrelatorConfig holds TOML-serialisable configuration for the
// event correlator pipeline stage.
type CorrelatorConfig struct {
	// Enabled controls whether correlation filtering is active.
	Enabled bool `toml:"enabled"`
	// WindowDuration is the sliding window expressed as a duration string
	// (e.g. "30s", "1m").
	WindowDuration string `toml:"window_duration"`
	// MinOccurrences is the minimum number of times a port must be
	// observed within the window before an alert is raised.
	MinOccurrences int `toml:"min_occurrences"`
}

// DefaultCorrelatorConfig returns the default correlator configuration.
func DefaultCorrelatorConfig() CorrelatorConfig {
	return CorrelatorConfig{
		Enabled:        true,
		WindowDuration: "30s",
		MinOccurrences: 2,
	}
}

// Validate returns an error if the configuration is invalid.
func (c CorrelatorConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	d, err := time.ParseDuration(c.WindowDuration)
	if err != nil {
		return fmt.Errorf("correlator: invalid window_duration %q: %w", c.WindowDuration, err)
	}
	if d <= 0 {
		return errorf("correlator: window_duration must be positive")
	}
	if c.MinOccurrences < 1 {
		return errorf("correlator: min_occurrences must be at least 1")
	}
	return nil
}

// WindowDurationParsed returns the WindowDuration as a time.Duration.
func (c CorrelatorConfig) WindowDurationParsed() (time.Duration, error) {
	return time.ParseDuration(c.WindowDuration)
}
