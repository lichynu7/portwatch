// Package config provides loading and validation of portwatch daemon
// configuration files.
//
// Configuration is expressed as YAML. Unknown fields are rejected to catch
// typos early. A call to Load merges the file over sensible defaults so that
// minimal config files remain valid:
//
//	# Minimum valid configuration file
//	interval: 60s
//	snapshot_dir: /var/lib/portwatch
//
// The Default function returns the same baseline without reading any file,
// which is useful for tests and for running portwatch without a config file.
//
// Validation rules:
//   - interval must be a positive duration (e.g. "30s", "5m")
//   - snapshot_dir must be a non-empty path
//   - Unknown top-level keys are rejected to surface typos early
//
// Example usage:
//
//	cfg, err := config.Load("/etc/portwatch/config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
package config
