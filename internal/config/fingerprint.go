package config

// FingerprintConfig controls whether port-set fingerprinting is enabled.
// When enabled, the daemon computes a hash of the current port set each
// scan cycle and skips alert dispatch when the set has not changed.
type FingerprintConfig struct {
	// Enabled turns fingerprint-based change detection on or off.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// LogOnMatch emits a debug-level log line when the fingerprint
	// matches the previous cycle (no change detected).
	LogOnMatch bool `toml:"log_on_match" yaml:"log_on_match"`
}

// DefaultFingerprintConfig returns the recommended defaults.
func DefaultFingerprintConfig() FingerprintConfig {
	return FingerprintConfig{
		Enabled:    true,
		LogOnMatch: false,
	}
}

// Validate checks that the FingerprintConfig is self-consistent.
func (c FingerprintConfig) Validate() error {
	// No numeric bounds to check; any bool combination is valid.
	return nil
}
