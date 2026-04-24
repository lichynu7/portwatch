// Package config loads and validates portwatch configuration.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Interval    time.Duration `yaml:"interval"`
	SnapshotDir string        `yaml:"snapshot_dir"`
	LogLevel    string        `yaml:"log_level"`
	Webhook     WebhookConfig `yaml:"webhook"`
	Email       EmailConfig   `yaml:"email"`
}

// WebhookConfig holds webhook alert settings.
type WebhookConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// EmailConfig holds SMTP alert settings.
type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Interval:    30 * time.Second,
		SnapshotDir: "/var/lib/portwatch",
		LogLevel:    "info",
		Webhook: WebhookConfig{
			Timeout: 5 * time.Second,
		},
	}
}

// Load reads a YAML config file from path, merging over defaults.
// Returns an error if the file cannot be read or the interval is invalid.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("config: read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse %s: %w", path, err)
	}
	if cfg.Interval <= 0 {
		return cfg, fmt.Errorf("config: interval must be positive")
	}
	return cfg, nil
}
