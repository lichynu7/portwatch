package ports

import (
	"sync"
	"time"
)

// CorrelatorConfig holds configuration for the event correlator.
type CorrelatorConfig struct {
	// WindowDuration is how long to collect events before correlating.
	WindowDuration time.Duration
	// MinOccurrences is the minimum number of times a port must appear
	// within the window to be considered a correlated (persistent) event.
	MinOccurrences int
}

// DefaultCorrelatorConfig returns a sensible default configuration.
func DefaultCorrelatorConfig() CorrelatorConfig {
	return CorrelatorConfig{
		WindowDuration: 30 * time.Second,
		MinOccurrences: 2,
	}
}

// Correlator tracks port appearances over a sliding window and emits
// only ports that have been seen at least MinOccurrences times, reducing
// noise from transient listeners.
type Correlator struct {
	cfg    CorrelatorConfig
	mu     sync.Mutex
	events map[string][]time.Time
	now    func() time.Time
}

// NewCorrelator creates a Correlator with the given config.
func NewCorrelator(cfg CorrelatorConfig) *Correlator {
	return &Correlator{
		cfg:    cfg,
		events: make(map[string][]time.Time),
		now:    time.Now,
	}
}

// Observe records a port observation and returns true if the port has
// been seen at least MinOccurrences times within the WindowDuration.
func (c *Correlator) Observe(p Port) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := portKey(p)
	cutoff := c.now().Add(-c.cfg.WindowDuration)

	// Prune old observations outside the window.
	filtered := c.events[key][:0]
	for _, t := range c.events[key] {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, c.now())
	c.events[key] = filtered

	return len(filtered) >= c.cfg.MinOccurrences
}

// Reset clears all tracked observations for a given port key.
func (c *Correlator) Reset(p Port) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.events, portKey(p))
}

// Purge removes all observation windows that have fully expired.
func (c *Correlator) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := c.now().Add(-c.cfg.WindowDuration)
	for key, times := range c.events {
		allExpired := true
		for _, t := range times {
			if t.After(cutoff) {
				allExpired = false
				break
			}
		}
		if allExpired {
			delete(c.events, key)
		}
	}
}
