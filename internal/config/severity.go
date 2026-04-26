package config

import "fmt"

// SeverityConfig mirrors ports.SeverityConfig for TOML/JSON unmarshalling
// and provides validation logic.
type SeverityConfig struct {
	PrivilegedMax     uint16   `toml:"privileged_max"     json:"privileged_max"`
	EphemeralMin      uint16   `toml:"ephemeral_min"      json:"ephemeral_min"`
	CriticalProcesses []string `toml:"critical_processes" json:"critical_processes"`
}

// DefaultSeverityConfig returns the default severity classification thresholds.
func DefaultSeverityConfig() SeverityConfig {
	return SeverityConfig{
		PrivilegedMax:     1023,
		EphemeralMin:      32768,
		CriticalProcesses: []string{"nc", "ncat", "nmap", "socat"},
	}
}

// Validate checks that the SeverityConfig fields are logically consistent.
func (c SeverityConfig) Validate() error {
	if c.PrivilegedMax == 0 {
		return fmt.Errorf("severity: privileged_max must be greater than 0")
	}
	if c.EphemeralMin == 0 {
		return fmt.Errorf("severity: ephemeral_min must be greater than 0")
	}
	if c.EphemeralMin <= c.PrivilegedMax {
		return fmt.Errorf(
			"severity: ephemeral_min (%d) must be greater than privileged_max (%d)",
			c.EphemeralMin, c.PrivilegedMax,
		)
	}
	return nil
}
