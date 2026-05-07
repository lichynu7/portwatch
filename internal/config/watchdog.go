package config

import "time"

// WatchdogConfig holds configuration for the pipeline stall watchdog.
type WatchdogConfig struct {
	// Enabled controls whether the watchdog runs.
	Enabled bool `toml:"enabled" yaml:"enabled"`
	// StallTimeout is how long without a scan tick before the stall callback fires.
	StallTimeout Duration `toml:"stall_timeout" yaml:"stall_timeout"`
}

// DefaultWatchdogConfig returns a production-ready watchdog configuration.
func DefaultWatchdogConfig() WatchdogConfig {
	return WatchdogConfig{
		Enabled:      true,
		StallTimeout: Duration(3 * time.Minute),
	}
}

// Validate returns an error if the configuration is invalid.
func (c WatchdogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if time.Duration(c.StallTimeout) <= 0 {
		return errorf("watchdog: stall_timeout must be positive, got %s", c.StallTimeout)
	}
	if time.Duration(c.StallTimeout) < 5*time.Second {
		return errorf("watchdog: stall_timeout %s is too short (minimum 5s)", c.StallTimeout)
	}
	return nil
}
