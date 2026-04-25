package config

// BaselineConfig holds settings for the port baseline feature.
// When enabled, portwatch loads a persisted baseline file and suppresses
// alerts for ports that an operator has explicitly acknowledged.
type BaselineConfig struct {
	// Enabled turns the baseline feature on or off.
	Enabled bool `toml:"enabled" json:"enabled"`

	// Path is the file system location of the baseline JSON file.
	// Defaults to /var/lib/portwatch/baseline.json.
	Path string `toml:"path" json:"path"`
}

// DefaultBaselineConfig returns a BaselineConfig with sensible defaults.
func DefaultBaselineConfig() BaselineConfig {
	return BaselineConfig{
		Enabled: false,
		Path:    "/var/lib/portwatch/baseline.json",
	}
}

// Validate checks that the BaselineConfig is coherent.
func (b BaselineConfig) Validate() error {
	if b.Enabled && b.Path == "" {
		return &ValidationError{Field: "baseline.path", Message: "path must not be empty when baseline is enabled"}
	}
	return nil
}

// ValidationError describes a configuration validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "config: " + e.Field + ": " + e.Message
}
