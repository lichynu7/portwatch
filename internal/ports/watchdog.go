package ports

import (
	"context"
	"sync"
	"time"
)

// WatchdogConfig holds configuration for the pipeline watchdog.
type WatchdogConfig struct {
	// Enabled controls whether the watchdog is active.
	Enabled bool
	// StallTimeout is how long without a scan tick before the watchdog fires.
	StallTimeout time.Duration
	// OnStall is called when a stall is detected.
	OnStall func(stalledFor time.Duration)
}

// DefaultWatchdogConfig returns a sensible default watchdog configuration.
func DefaultWatchdogConfig() WatchdogConfig {
	return WatchdogConfig{
		Enabled:      true,
		StallTimeout: 3 * time.Minute,
		OnStall:      nil,
	}
}

// Watchdog monitors the pipeline heartbeat and fires a callback when
// no tick has been received within StallTimeout.
type Watchdog struct {
	cfg     WatchdogConfig
	mu      sync.Mutex
	lastSeen time.Time
}

// NewWatchdog creates a Watchdog from cfg. Returns nil if disabled.
func NewWatchdog(cfg WatchdogConfig) *Watchdog {
	if !cfg.Enabled {
		return nil
	}
	return &Watchdog{
		cfg:      cfg,
		lastSeen: time.Now(),
	}
}

// Tick records that the pipeline is alive. Call this after each scan cycle.
func (w *Watchdog) Tick() {
	w.mu.Lock()
	w.lastSeen = time.Now()
	w.mu.Unlock()
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.StallTimeout / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.mu.Lock()
			stalledFor := time.Since(w.lastSeen)
			w.mu.Unlock()
			if stalledFor >= w.cfg.StallTimeout && w.cfg.OnStall != nil {
				w.cfg.OnStall(stalledFor)
			}
		}
	}
}
