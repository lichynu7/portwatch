package config

import "fmt"

// MetadataConfig controls the per-port metadata annotation store.
type MetadataConfig struct {
	// Enabled controls whether metadata collection is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// MaxKeysPerPort is the maximum number of metadata keys stored per port.
	// Defaults to 16.
	MaxKeysPerPort int `toml:"max_keys_per_port" yaml:"max_keys_per_port"`
}

// DefaultMetadataConfig returns a MetadataConfig with sensible defaults.
func DefaultMetadataConfig() MetadataConfig {
	return MetadataConfig{
		Enabled:        true,
		MaxKeysPerPort: 16,
	}
}

// Validate checks that the MetadataConfig is self-consistent.
func (c MetadataConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxKeysPerPort <= 0 {
		return fmt.Errorf("metadata: max_keys_per_port must be > 0, got %d", c.MaxKeysPerPort)
	}
	if c.MaxKeysPerPort > 256 {
		return fmt.Errorf("metadata: max_keys_per_port must be <= 256, got %d", c.MaxKeysPerPort)
	}
	return nil
}
