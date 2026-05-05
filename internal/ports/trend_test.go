package ports

import (
	"testing"
	"time"
)

func TestTrendStableWithSingleSample(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 2)
	now := time.Now()
	td.Record("tcp:8080", 3, now)

	res := td.Evaluate("tcp:8080")
	if res.Direction != TrendStable {
		t.Errorf("expected stable with single sample, got %s", res.Direction)
	}
	if res.Delta != 0 {
		t.Errorf("expected delta 0, got %d", res.Delta)
	}
}

func TestTrendRising(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 2)
	now := time.Now()
	td.Record("tcp:9000", 1, now.Add(-2*time.Minute))
	td.Record("tcp:9000", 5, now)

	res := td.Evaluate("tcp:9000")
	if res.Direction != TrendRising {
		t.Errorf("expected rising, got %s", res.Direction)
	}
	if res.Delta != 4 {
		t.Errorf("expected delta 4, got %d", res.Delta)
	}
}

func TestTrendFalling(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 2)
	now := time.Now()
	td.Record("tcp:443", 10, now.Add(-3*time.Minute))
	td.Record("tcp:443", 2, now)

	res := td.Evaluate("tcp:443")
	if res.Direction != TrendFalling {
		t.Errorf("expected falling, got %s", res.Direction)
	}
	if res.Delta != -8 {
		t.Errorf("expected delta -8, got %d", res.Delta)
	}
}

func TestTrendStableSmallDelta(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 3)
	now := time.Now()
	td.Record("udp:53", 4, now.Add(-1*time.Minute))
	td.Record("udp:53", 5, now)

	res := td.Evaluate("udp:53")
	if res.Direction != TrendStable {
		t.Errorf("expected stable for delta below minDelta, got %s", res.Direction)
	}
}

func TestTrendExpiredSamplesEvicted(t *testing.T) {
	td := NewTrendDetector(1*time.Minute, 2)
	now := time.Now()
	// old sample outside window
	td.Record("tcp:22", 100, now.Add(-5*time.Minute))
	// recent samples
	td.Record("tcp:22", 3, now.Add(-30*time.Second))
	td.Record("tcp:22", 4, now)

	res := td.Evaluate("tcp:22")
	// delta should be computed only over in-window samples: 4-3=1, stable
	if res.Direction != TrendStable {
		t.Errorf("expected stable after eviction, got %s", res.Direction)
	}
}

func TestTrendIndependentKeys(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 2)
	now := time.Now()
	td.Record("tcp:80", 1, now.Add(-1*time.Minute))
	td.Record("tcp:80", 10, now)
	td.Record("tcp:443", 5, now.Add(-1*time.Minute))
	td.Record("tcp:443", 5, now)

	if td.Evaluate("tcp:80").Direction != TrendRising {
		t.Error("tcp:80 should be rising")
	}
	if td.Evaluate("tcp:443").Direction != TrendStable {
		t.Error("tcp:443 should be stable")
	}
}

func TestTrendReset(t *testing.T) {
	td := NewTrendDetector(5*time.Minute, 2)
	now := time.Now()
	td.Record("tcp:8443", 1, now.Add(-1*time.Minute))
	td.Record("tcp:8443", 9, now)
	td.Reset("tcp:8443")

	res := td.Evaluate("tcp:8443")
	if res.Direction != TrendStable {
		t.Errorf("expected stable after reset, got %s", res.Direction)
	}
}
