package ports

import (
	"testing"
	"time"
)

func TestAnomalyFirstObservationBelowMinHits(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 3)
	if det.Observe("tcp:8080") {
		t.Fatal("expected false on first observation with minHits=3")
	}
}

func TestAnomalyReachesMinHits(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 3)
	det.Observe("tcp:8080")
	det.Observe("tcp:8080")
	if !det.Observe("tcp:8080") {
		t.Fatal("expected true after reaching minHits")
	}
}

func TestAnomalyWindowExpiry(t *testing.T) {
	now := time.Now()
	det := NewAnomalyDetector(time.Minute, 2)
	det.now = func() time.Time { return now }

	det.Observe("tcp:9000")

	// Advance time beyond the window.
	det.now = func() time.Time { return now.Add(2 * time.Minute) }

	// Should reset; first observation in new window.
	if det.Observe("tcp:9000") {
		t.Fatal("expected false after window expiry with minHits=2")
	}
}

func TestAnomalyMinHitsOne(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 1)
	if !det.Observe("tcp:22") {
		t.Fatal("expected true immediately when minHits=1")
	}
}

func TestAnomalyIndependentKeys(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2)
	det.Observe("tcp:80")
	det.Observe("tcp:80")

	// A different key should not be affected.
	if det.Observe("tcp:443") {
		t.Fatal("tcp:443 should not be flagged after one observation")
	}
}

func TestAnomalyReset(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2)
	det.Observe("tcp:8443")
	det.Reset("tcp:8443")

	// After reset the counter starts from zero.
	if det.Observe("tcp:8443") {
		t.Fatal("expected false after reset with minHits=2")
	}
}

func TestAnomalyPurgeRemovesStale(t *testing.T) {
	now := time.Now()
	det := NewAnomalyDetector(time.Minute, 2)
	det.now = func() time.Time { return now }
	det.Observe("tcp:1234")

	// Advance beyond window so the record is stale.
	det.now = func() time.Time { return now.Add(5 * time.Minute) }
	det.Purge()

	// After purge the record is gone; next Observe starts fresh.
	if det.Observe("tcp:1234") {
		t.Fatal("expected false after purge with minHits=2")
	}
}
