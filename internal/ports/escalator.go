package ports

import (
	"sync"
	"time"
)

// EscalatorConfig controls how repeated alerts escalate in severity.
type EscalatorConfig struct {
	// Window is the duration over which repeated occurrences are counted.
	Window time.Duration
	// Threshold is the number of occurrences required to escalate.
	Threshold int
	// Enabled controls whether escalation is active.
	Enabled bool
}

// DefaultEscalatorConfig returns a sensible default configuration.
func DefaultEscalatorConfig() EscalatorConfig {
	return EscalatorConfig{
		Window:    5 * time.Minute,
		Threshold: 3,
		Enabled:   true,
	}
}

// escalatorEntry tracks occurrences for a single port key.
type escalatorEntry struct {
	count     int
	firstSeen time.Time
}

// Escalator promotes alert severity when a port is seen repeatedly
// within a configured time window.
type Escalator struct {
	cfg     EscalatorConfig
	mu      sync.Mutex
	entries map[string]*escalatorEntry
	now     func() time.Time
}

// NewEscalator creates a new Escalator with the given configuration.
func NewEscalator(cfg EscalatorConfig) *Escalator {
	return &Escalator{
		cfg:     cfg,
		entries: make(map[string]*escalatorEntry),
		now:     time.Now,
	}
}

// Observe records an occurrence for the given key and returns true when
// the threshold has been reached, indicating the alert should escalate.
func (e *Escalator) Observe(key string) bool {
	if !e.cfg.Enabled {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	ent, ok := e.entries[key]
	if !ok || now.Sub(ent.firstSeen) > e.cfg.Window {
		e.entries[key] = &escalatorEntry{count: 1, firstSeen: now}
		return false
	}

	ent.count++
	return ent.count >= e.cfg.Threshold
}

// Reset clears the occurrence history for a given key.
func (e *Escalator) Reset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, key)
}

// Purge removes all entries whose window has expired.
func (e *Escalator) Purge() {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	for k, ent := range e.entries {
		if now.Sub(ent.firstSeen) > e.cfg.Window {
			delete(e.entries, k)
		}
	}
}
