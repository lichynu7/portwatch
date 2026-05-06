package config

import "fmt"

// ChangelogConfig controls the in-memory port change history.
type ChangelogConfig struct {
	// Enabled toggles changelog collection entirely.
	Enabled bool `toml:"enabled" yaml:"enabled"`
	// MaxEvents is the maximum number of events retained in memory.
	MaxEvents int `toml:"max_events" yaml:"max_events"`
}

// DefaultChangelogConfig returns a sensible out-of-the-box configuration.
func DefaultChangelogConfig() ChangelogConfig {
	return ChangelogConfig{
		Enabled:   true,
		MaxEvents: 256,
	}
}

// Validate returns an error if the configuration is invalid.
func (c ChangelogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxEvents <= 0 {
		return fmt.Errorf("changelog: max_events must be > 0, got %d", c.MaxEvents)
	}
	if c.MaxEvents > 10_000 {
		return fmt.Errorf("changelog: max_events %d exceeds maximum of 10000", c.MaxEvents)
	}
	return nil
}
