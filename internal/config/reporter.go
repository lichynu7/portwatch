package config

import "fmt"

// ReporterConfig holds configuration for the port alert reporter.
type ReporterConfig struct {
	// Format is the output format: "table" or "json".
	Format string `toml:"format" yaml:"format"`

	// OutputFile is an optional path to write reports to.
	// If empty, output goes to stdout.
	OutputFile string `toml:"output_file" yaml:"output_file"`

	// IncludeFields lists which fields to include in the report.
	// An empty slice means all fields are included.
	IncludeFields []string `toml:"include_fields" yaml:"include_fields"`
}

var validFormats = map[string]bool{
	"table": true,
	"json":  true,
}

// DefaultReporterConfig returns sensible defaults for the reporter.
func DefaultReporterConfig() ReporterConfig {
	return ReporterConfig{
		Format:        "table",
		OutputFile:    "",
		IncludeFields: []string{},
	}
}

// Validate checks that the ReporterConfig is well-formed.
func (c ReporterConfig) Validate() error {
	if c.Format == "" {
		return errorf("reporter", "format must not be empty")
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("reporter: unsupported format %q; valid options: table, json", c.Format)
	}
	return nil
}
