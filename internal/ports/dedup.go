package ports

import (
	"sync"
	"time"
)

// DedupConfig holds configuration for the deduplication window.
type DedupConfig struct {
	// WindowSize is how long a seen port key is remembered.
	WindowSize time.Duration
}

// DefaultDedupConfig returns a sensible default dedup configuration.
func DefaultDedupConfig() DedupConfig {
	return DedupConfig{
		WindowSize: 5 * time.Minute,
	}
}

// Deduplicator suppresses repeated alerts for the same port within a time window.
type Deduplicator struct {
	mu     sync.Mutex
	seen   map[string]time.Time
	window time.Duration
	now    func() time.Time
}

// NewDeduplicator creates a Deduplicator with the given window duration.
func NewDeduplicator(cfg DedupConfig) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: cfg.WindowSize,
		now:    time.Now,
	}
}

// IsDuplicate returns true if the key was seen within the dedup window.
// If not a duplicate, it records the key as seen.
func (d *Deduplicator) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if t, ok := d.seen[key]; ok && now.Sub(t) < d.window {
		return true
	}
	d.seen[key] = now
	return false
}

// Purge removes all entries that have expired beyond the window.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Reset clears all seen entries.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
