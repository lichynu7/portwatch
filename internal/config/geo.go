package config

import "fmt"

// GeoConfig controls the geo-classification enrichment stage.
type GeoConfig struct {
	// Enabled turns geo classification on or off.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// FilterLabels, when non-empty, restricts alerts to ports whose geo label
	// matches one of the listed values (e.g. "public", "private").
	FilterLabels []string `toml:"filter_labels" yaml:"filter_labels"`
}

// DefaultGeoConfig returns a GeoConfig with sensible defaults.
func DefaultGeoConfig() GeoConfig {
	return GeoConfig{
		Enabled:      true,
		FilterLabels: []string{},
	}
}

// Validate checks the GeoConfig for invalid combinations.
func (g GeoConfig) Validate() error {
	valid := map[string]bool{
		"loopback":   true,
		"link-local": true,
		"private":    true,
		"public":     true,
		"invalid":    true,
	}
	for _, label := range g.FilterLabels {
		if !valid[label] {
			return fmt.Errorf("geo: unknown filter label %q", label)
		}
	}
	return nil
}
