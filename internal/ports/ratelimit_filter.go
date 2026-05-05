package ports

import (
	"sync"
	"time"
)

// PortEvent represents a single alert event for a port.
type PortEvent struct {
	Key       string
	Occurred  time.Time
	Suppressed bool
}

// AlertRateLimitFilter drops repeated alerts for the same port key within a
// sliding window. Unlike the deduplicator it tracks per-key hit counts and
// only suppresses once a burst ceiling is reached.
type AlertRateLimitFilter struct {
	cfg    RateLimitFilterConfig
	mu     sync.Mutex
	bucket map[string]*rlBucket
}

type rlBucket struct {
	hits      int
	windowEnd time.Time
}

// RateLimitFilterConfig holds tuning knobs for AlertRateLimitFilter.
type RateLimitFilterConfig struct {
	// Window is the rolling duration in which hits are counted.
	Window time.Duration
	// MaxHits is the number of alerts allowed per window before suppression.
	MaxHits int
}

// DefaultRateLimitFilterConfig returns sensible defaults.
func DefaultRateLimitFilterConfig() RateLimitFilterConfig {
	return RateLimitFilterConfig{
		Window:  30 * time.Second,
		MaxHits: 3,
	}
}

// NewAlertRateLimitFilter constructs a filter with the given config.
func NewAlertRateLimitFilter(cfg RateLimitFilterConfig) *AlertRateLimitFilter {
	return &AlertRateLimitFilter{
		cfg:    cfg,
		bucket: make(map[string]*rlBucket),
	}
}

// Allow returns true when the event should be forwarded and false when it
// should be suppressed.
func (f *AlertRateLimitFilter) Allow(key string, now time.Time) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	b, ok := f.bucket[key]
	if !ok || now.After(b.windowEnd) {
		f.bucket[key] = &rlBucket{hits: 1, windowEnd: now.Add(f.cfg.Window)}
		return true
	}
	b.hits++
	return b.hits <= f.cfg.MaxHits
}

// Reset clears state for a specific key.
func (f *AlertRateLimitFilter) Reset(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.bucket, key)
}

// ApplyRateLimitFilter wraps a slice of Port values, suppressing those whose
// key exceeds the configured burst ceiling within the current window.
func ApplyRateLimitFilter(ports []Port, f *AlertRateLimitFilter, now time.Time) []Port {
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		if f.Allow(portKey(p), now) {
			out = append(out, p)
		}
	}
	return out
}
