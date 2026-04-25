package config

import "time"

// ThrottleConfig holds user-facing throttle settings loaded from the config
// file. It mirrors ports.ThrottleConfig but uses plain types for TOML/YAML
// unmarshalling.
type ThrottleConfig struct {
	// WindowSeconds is the cooldown period in seconds between repeated alerts
	// for the same port. Defaults to 300 (5 minutes).
	WindowSeconds int `toml:"window_seconds" yaml:"window_seconds"`
	// MaxBurst is the number of alerts permitted before throttling engages
	// within a window. Defaults to 1.
	MaxBurst int `toml:"max_burst" yaml:"max_burst"`
}

// DefaultThrottleConfig returns conservative defaults suitable for most
// production deployments.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		WindowSeconds: 300,
		MaxBurst:      1,
	}
}

// Window converts WindowSeconds to a time.Duration.
func (c ThrottleConfig) Window() time.Duration {
	if c.WindowSeconds <= 0 {
		return 5 * time.Minute
	}
	return time.Duration(c.WindowSeconds) * time.Second
}

// Validate returns an error string if the configuration is invalid, or an
// empty string when the config is acceptable.
func (c ThrottleConfig) Validate() string {
	if c.WindowSeconds < 0 {
		return "throttle window_seconds must be >= 0"
	}
	if c.MaxBurst < 1 {
		return "throttle max_burst must be >= 1"
	}
	return ""
}
