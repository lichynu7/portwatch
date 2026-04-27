package config

import "fmt"

// SeverityFilterConfig holds configuration for the severity-based port filter.
type SeverityFilterConfig struct {
	// MinLevel is the minimum severity level a port must reach to be reported.
	// Accepted values: "info", "warning", "critical".
	MinLevel string `toml:"min_level" yaml:"min_level"`
}

// DefaultSeverityFilterConfig returns a SeverityFilterConfig with sensible
// defaults (report everything at "info" level and above).
func DefaultSeverityFilterConfig() SeverityFilterConfig {
	return SeverityFilterConfig{
		MinLevel: "info",
	}
}

// Validate checks that MinLevel is a recognised severity string.
func (c SeverityFilterConfig) Validate() error {
	switch c.MinLevel {
	case "info", "warning", "critical":
		return nil
	default:
		return fmt.Errorf("severity_filter: unknown min_level %q (want info|warning|critical)", c.MinLevel)
	}
}
