package ports

import (
	"sync"
	"time"
)

// TrendDirection indicates whether port activity is increasing, stable, or decreasing.
type TrendDirection string

const (
	TrendRising  TrendDirection = "rising"
	TrendStable  TrendDirection = "stable"
	TrendFalling TrendDirection = "falling"
)

// TrendSample records a count observation at a point in time.
type TrendSample struct {
	At    time.Time
	Count int
}

// TrendResult holds the computed trend for a given key.
type TrendResult struct {
	Key       string
	Direction TrendDirection
	Delta     int // difference between newest and oldest sample in window
}

// TrendDetector tracks port observation counts over a sliding window and
// classifies whether activity is rising, stable, or falling.
type TrendDetector struct {
	mu       sync.Mutex
	window   time.Duration
	minDelta int
	samples  map[string][]TrendSample
}

// NewTrendDetector creates a TrendDetector with the given sliding window and
// minimum absolute delta required to classify a trend as non-stable.
func NewTrendDetector(window time.Duration, minDelta int) *TrendDetector {
	if window <= 0 {
		window = 5 * time.Minute
	}
	if minDelta <= 0 {
		minDelta = 2
	}
	return &TrendDetector{
		window:   window,
		minDelta: minDelta,
		samples:  make(map[string][]TrendSample),
	}
}

// Record adds a new count observation for the given key at now.
func (t *TrendDetector) Record(key string, count int, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.window)
	existing := t.samples[key]
	filtered := existing[:0]
	for _, s := range existing {
		if s.At.After(cutoff) {
			filtered = append(filtered, s)
		}
	}
	filtered = append(filtered, TrendSample{At: now, Count: count})
	t.samples[key] = filtered
}

// Evaluate returns the current trend for the given key based on recorded samples.
func (t *TrendDetector) Evaluate(key string) TrendResult {
	t.mu.Lock()
	defer t.mu.Unlock()

	samples := t.samples[key]
	if len(samples) < 2 {
		return TrendResult{Key: key, Direction: TrendStable, Delta: 0}
	}

	first := samples[0].Count
	last := samples[len(samples)-1].Count
	delta := last - first

	dir := TrendStable
	if delta >= t.minDelta {
		dir = TrendRising
	} else if delta <= -t.minDelta {
		dir = TrendFalling
	}

	return TrendResult{Key: key, Direction: dir, Delta: delta}
}

// Reset clears all recorded samples for the given key.
func (t *TrendDetector) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.samples, key)
}
