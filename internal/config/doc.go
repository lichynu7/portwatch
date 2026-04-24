// Package config provides configuration loading and validation for portwatch.
//
// Configuration can be supplied via a TOML file. When no file is present the
// package returns safe defaults so that portwatch is immediately usable without
// any setup.
//
// # File format
//
// A minimal configuration file looks like:
//
//	[portwatch]
//	interval       = "30s"
//	snapshot_path  = "/var/lib/portwatch/snapshot.json"
//
//	[alert.slack]
//	webhook_url = "https://hooks.slack.com/services/…"
//
//	[alert.email]
//	smtp_host   = "smtp.example.com"
//	smtp_port   = 587
//	recipients  = ["ops@example.com"]
//
//	[alert.webhook]
//	url = "https://example.com/portwatch"
//
// # Alert handlers
//
// Each alert section is optional. Omitting a section disables that handler.
// Multiple handlers can be active simultaneously; portwatch fans out every
// alert to all registered handlers.
package config
