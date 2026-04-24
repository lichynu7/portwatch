package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	Interval    time.Duration `yaml:"interval"`
	SnapshotDir string        `yaml:"snapshot_dir"`
	AllowedPorts []int        `yaml:"allowed_ports"`
	Alerts      AlertConfig   `yaml:"alerts"`
}

// AlertConfig configures alert delivery channels.
type AlertConfig struct {
	LogFile string `yaml:"log_file"`
	Stdout  bool   `yaml:"stdout"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Interval:    30 * time.Second,
		SnapshotDir: "/var/lib/portwatch",
		Alerts: AlertConfig{
			Stdout: true,
		},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %s: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: invalid: %w", err)
	}

	return cfg, nil
}

// validate checks that required fields are sane.
func (c *Config) validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1s, got %s", c.Interval)
	}
	if c.SnapshotDir == "" {
		return fmt.Errorf("snapshot_dir must not be empty")
	}
	return nil
}
