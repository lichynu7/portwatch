// Package config holds configuration types and loaders for portwatch.
package config

import "errors"

// EnricherConfig controls port metadata enrichment behaviour.
type EnricherConfig struct {
	// ResolveServices enables mapping port numbers to well-known service names.
	ResolveServices bool `toml:"resolve_services" yaml:"resolve_services"`

	// LookupProcesses enables enriching ports with owning process information
	// by inspecting /proc. Requires appropriate OS permissions.
	LookupProcesses bool `toml:"lookup_processes" yaml:"lookup_processes"`
}

// DefaultEnricherConfig returns a sensible out-of-the-box enricher config.
func DefaultEnricherConfig() EnricherConfig {
	return EnricherConfig{
		ResolveServices: true,
		LookupProcesses: true,
	}
}

// Validate checks that the EnricherConfig is self-consistent.
func (c EnricherConfig) Validate() error {
	// Currently no invalid combinations; reserved for future constraints.
	if !c.ResolveServices && !c.LookupProcesses {
		// Both disabled is valid — enricher simply passes through raw ports.
		return nil
	}
	if c.LookupProcesses && !isLinux() {
		return errors.New("enricher: process lookup is only supported on Linux")
	}
	return nil
}

// isLinux reports whether the current OS is Linux.
// Defined as a variable so tests can override it.
var isLinux = func() bool {
	return goosIsLinux()
}
