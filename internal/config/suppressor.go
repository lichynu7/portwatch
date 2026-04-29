package config

import (
	"fmt"

	"github.com/user/portwatch/internal/ports"
)

// SuppressorConfig holds the TOML-decoded configuration for the alert suppressor.
type SuppressorConfig struct {
	// QuietHoursStart is the hour (0-23) when suppression begins. -1 disables.
	QuietHoursStart int `toml:"quiet_hours_start"`
	// QuietHoursEnd is the hour (0-23) when suppression ends. -1 disables.
	QuietHoursEnd int `toml:"quiet_hours_end"`
	// MinSeverity is the severity label that bypasses quiet hours (e.g. "critical").
	MinSeverity string `toml:"min_severity"`
}

// DefaultSuppressorConfig returns a SuppressorConfig with quiet hours disabled.
func DefaultSuppressorConfig() SuppressorConfig {
	return SuppressorConfig{
		QuietHoursStart: -1,
		QuietHoursEnd:   -1,
		MinSeverity:     "critical",
	}
}

// Validate checks the SuppressorConfig for logical errors.
func (c SuppressorConfig) Validate() error {
	if c.QuietHoursStart == -1 && c.QuietHoursEnd == -1 {
		return nil // disabled
	}
	if c.QuietHoursStart < 0 || c.QuietHoursStart > 23 {
		return errorf("suppressor: quiet_hours_start must be 0-23 or -1, got %d", c.QuietHoursStart)
	}
	if c.QuietHoursEnd < 0 || c.QuietHoursEnd > 23 {
		return errorf("suppressor: quiet_hours_end must be 0-23 or -1, got %d", c.QuietHoursEnd)
	}
	if _, err := parseSeverity(c.MinSeverity); err != nil {
		return fmt.Errorf("suppressor: %w", err)
	}
	return nil
}

// ToPortsSuppressorConfig converts this config into the ports-layer struct.
func (c SuppressorConfig) ToPortsSuppressorConfig() ports.SuppressorConfig {
	minSev, _ := parseSeverity(c.MinSeverity)
	return ports.SuppressorConfig{
		QuietHoursStart: c.QuietHoursStart,
		QuietHoursEnd:   c.QuietHoursEnd,
		MinSeverity:     minSev,
	}
}

func parseSeverity(s string) (ports.Severity, error) {
	switch s {
	case "info":
		return ports.SeverityInfo, nil
	case "warning":
		return ports.SeverityWarning, nil
	case "critical":
		return ports.SeverityCritical, nil
	default:
		return 0, fmt.Errorf("unknown severity %q: must be info, warning, or critical", s)
	}
}
