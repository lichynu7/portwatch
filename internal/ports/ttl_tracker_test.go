package ports

import (
	"testing"
	"time"
)

func TestTTLTrackerObserveCreatesEntry(t *testing.T) {
	tr := NewTTLTracker(5 * time.Minute)
	e := tr.Observe("tcp:8080")
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tr.Len())
	}
}

func TestTTLTrackerObserveIncrementsCount(t *testing.T) {
	tr := NewTTLTracker(5 * time.Minute)
	tr.Observe("tcp:9090")
	e := tr.Observe("tcp:9090")
	if e.Count != 2 {
		t.Fatalf("expected count 2, got %d", e.Count)
	}
}

func TestTTLTrackerAgeUnknownKeyIsZero(t *testing.T) {
	tr := NewTTLTracker(5 * time.Minute)
	if d := tr.Age("tcp:1234"); d != 0 {
		t.Fatalf("expected 0 age for unknown key, got %v", d)
	}
}

func TestTTLTrackerAgeGrows(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := NewTTLTracker(5 * time.Minute)
	tr.now = func() time.Time { return base }
	tr.Observe("tcp:8080")

	tr.now = func() time.Time { return base.Add(30 * time.Second) }
	age := tr.Age("tcp:8080")
	if age != 30*time.Second {
		t.Fatalf("expected 30s age, got %v", age)
	}
}

func TestTTLTrackerEvictRemovesStale(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := NewTTLTracker(1 * time.Minute)
	tr.now = func() time.Time { return base }
	tr.Observe("tcp:8080")
	tr.Observe("tcp:9090")

	// Advance past TTL for both
	tr.now = func() time.Time { return base.Add(2 * time.Minute) }
	removed := tr.Evict()
	if removed != 2 {
		t.Fatalf("expected 2 evictions, got %d", removed)
	}
	if tr.Len() != 0 {
		t.Fatalf("expected empty tracker, got %d", tr.Len())
	}
}

func TestTTLTrackerEvictKeepsFresh(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := NewTTLTracker(1 * time.Minute)
	tr.now = func() time.Time { return base }
	tr.Observe("tcp:8080")

	// Re-observe just before eviction window
	tr.now = func() time.Time { return base.Add(59 * time.Second) }
	tr.Observe("tcp:8080")

	tr.now = func() time.Time { return base.Add(90 * time.Second) }
	removed := tr.Evict()
	if removed != 0 {
		t.Fatalf("expected 0 evictions, got %d", removed)
	}
}

func TestTTLTrackerIndependentKeys(t *testing.T) {
	tr := NewTTLTracker(5 * time.Minute)
	tr.Observe("tcp:80")
	tr.Observe("tcp:80")
	tr.Observe("udp:53")

	e80 := tr.Observe("tcp:80")
	e53 := tr.Observe("udp:53")
	if e80.Count != 3 {
		t.Fatalf("tcp:80 count want 3, got %d", e80.Count)
	}
	if e53.Count != 2 {
		t.Fatalf("udp:53 count want 2, got %d", e53.Count)
	}
}
