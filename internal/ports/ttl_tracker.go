package ports

import (
	"sync"
	"time"
)

// TTLEntry holds the first-seen time and last-seen time for a port key.
type TTLEntry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// TTLTracker tracks how long each port has been continuously observed,
// evicting entries that have not been seen within the TTL window.
type TTLTracker struct {
	mu      sync.Mutex
	entries map[string]*TTLEntry
	ttl     time.Duration
	now     func() time.Time
}

// NewTTLTracker creates a TTLTracker with the given TTL duration.
func NewTTLTracker(ttl time.Duration) *TTLTracker {
	return &TTLTracker{
		entries: make(map[string]*TTLEntry),
		ttl:     ttl,
		now:     time.Now,
	}
}

// Observe records an observation for the given key and returns the entry.
func (t *TTLTracker) Observe(key string) TTLEntry {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e, ok := t.entries[key]
	if !ok {
		e = &TTLEntry{FirstSeen: now}
		t.entries[key] = e
	}
	e.LastSeen = now
	e.Count++
	return *e
}

// Evict removes entries whose LastSeen is older than the TTL.
func (t *TTLTracker) Evict() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := t.now().Add(-t.ttl)
	removed := 0
	for k, e := range t.entries {
		if e.LastSeen.Before(cutoff) {
			delete(t.entries, k)
			removed++
		}
	}
	return removed
}

// Age returns how long a key has been tracked. Returns 0 if unknown.
func (t *TTLTracker) Age(key string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if e, ok := t.entries[key]; ok {
		return t.now().Sub(e.FirstSeen)
	}
	return 0
}

// Len returns the number of tracked entries.
func (t *TTLTracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
