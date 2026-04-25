package ports

import (
	"sync"
	"time"
)

// RateLimiter tracks how frequently a given port has triggered alerts
// and suppresses repeated alerts within a cooldown window.
type RateLimiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSeen map[uint16]time.Time
}

// NewRateLimiter creates a RateLimiter with the given cooldown duration.
// Alerts for a port are suppressed if the same port was seen within the
// cooldown window.
func NewRateLimiter(cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		cooldown: cooldown,
		lastSeen: make(map[uint16]time.Time),
	}
}

// Allow returns true if an alert for the given port should be delivered.
// It returns false if the port was already alerted within the cooldown window.
// When Allow returns true it records the current time for that port.
func (r *RateLimiter) Allow(port uint16) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if last, ok := r.lastSeen[port]; ok {
		if time.Since(last) < r.cooldown {
			return false
		}
	}
	r.lastSeen[port] = time.Now()
	return true
}

// Reset clears the rate-limit state for a specific port.
// Useful when a port disappears and later reappears.
func (r *RateLimiter) Reset(port uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.lastSeen, port)
}

// Purge removes all entries older than the cooldown window, keeping
// memory usage bounded during long daemon runs.
func (r *RateLimiter) Purge() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for port, last := range r.lastSeen {
		if now.Sub(last) >= r.cooldown {
			delete(r.lastSeen, port)
		}
	}
}
