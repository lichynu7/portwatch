package config

import "time"

// WatcherConfig holds tunable parameters for the port watcher.
type WatcherConfig struct {
	// Interval between consecutive port scans.
	Interval time.Duration `toml:"interval" yaml:"interval"`

	// ProcFS is the path to the proc filesystem (default: /proc).
	ProcFS string `toml:"proc_fs" yaml:"proc_fs"`

	// EmitClosed controls whether "closed" events are forwarded to alert handlers.
	EmitClosed bool `toml:"emit_closed" yaml:"emit_closed"`

	// BufferSize is the capacity of the internal event channel.
	BufferSize int `toml:"buffer_size" yaml:"buffer_size"`
}

// DefaultWatcherConfig returns a WatcherConfig populated with sensible defaults.
func DefaultWatcherConfig() WatcherConfig {
	return WatcherConfig{
		Interval:   5 * time.Second,
		ProcFS:     "/proc",
		EmitClosed: false,
		BufferSize: 64,
	}
}

// Validate returns an error if the WatcherConfig contains invalid values.
func (w WatcherConfig) Validate() error {
	if w.Interval <= 0 {
		return errorf("watcher.interval must be positive, got %s", w.Interval)
	}
	if w.BufferSize < 1 {
		return errorf("watcher.buffer_size must be at least 1, got %d", w.BufferSize)
	}
	if w.ProcFS == "" {
		return errorf("watcher.proc_fs must not be empty")
	}
	return nil
}
