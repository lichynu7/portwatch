package ports

import (
	"sync"
	"time"
)

// AnomalyRecord tracks how often a port has been seen across scan cycles.
type AnomalyRecord struct {
	FirstSeen  time.Time
	LastSeen   time.Time
	Occurrences int
}

// AnomalyDetector flags ports that appear suddenly and persist beyond a
// minimum observation window, distinguishing transient noise from real
// listeners.
type AnomalyDetector struct {
	mu      sync.Mutex
	records map[string]*AnomalyRecord
	window  time.Duration
	minHits int
	now     func() time.Time
}

// DefaultAnomalyWindow is the minimum duration a port must be seen before it
// is considered an anomaly rather than a transient connection.
const DefaultAnomalyWindow = 2 * time.Minute

// NewAnomalyDetector creates a detector that requires a port to appear at
// least minHits times within window before raising an anomaly.
func NewAnomalyDetector(window time.Duration, minHits int) *AnomalyDetector {
	if window <= 0 {
		window = DefaultAnomalyWindow
	}
	if minHits < 1 {
		minHits = 1
	}
	return &AnomalyDetector{
		records: make(map[string]*AnomalyRecord),
		window:  window,
		minHits: minHits,
		now:     time.Now,
	}
}

// Observe records a sighting of key (typically protocol:port) and returns
// true when the port has been seen enough times within the configured window
// to be flagged as an anomaly.
func (a *AnomalyDetector) Observe(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.now()
	rec, ok := a.records[key]
	if !ok {
		a.records[key] = &AnomalyRecord{
			FirstSeen:   now,
			LastSeen:    now,
			Occurrences: 1,
		}
		return a.minHits <= 1
	}

	// Reset if the previous observation is outside the window.
	if now.Sub(rec.FirstSeen) > a.window {
		rec.FirstSeen = now
		rec.Occurrences = 1
		rec.LastSeen = now
		return a.minHits <= 1
	}

	rec.Occurrences++
	rec.LastSeen = now
	return rec.Occurrences >= a.minHits
}

// Reset clears all observation history for key.
func (a *AnomalyDetector) Reset(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.records, key)
}

// Purge removes records whose last observation is older than the configured
// window, keeping memory bounded during long daemon runs.
func (a *AnomalyDetector) Purge() {
	a.mu.Lock()
	defer a.mu.Unlock()
	cutoff := a.now().Add(-a.window)
	for k, r := range a.records {
		if r.LastSeen.Before(cutoff) {
			delete(a.records, k)
		}
	}
}
