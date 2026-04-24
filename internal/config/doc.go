// Package config provides loading and validation of portwatch daemon
// configuration files.
//
// Configuration is expressed as YAML. Unknown fields are rejected to catch
// typos early. A call to Load merges the file over sensible defaults so that
// minimal config files remain valid:
//
//	interval: 60s
//	snapshot_dir: /var/lib/portwatch
//
// The Default function returns the same baseline without reading any file,
// which is useful for tests and for running portwatch without a config file.
package config
