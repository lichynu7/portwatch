package ports

import (
	"testing"
	"time"
)

func makeCorrelatorPort(port uint16) Port {
	return Port{
		LocalAddress: "0.0.0.0",
		LocalPort:    port,
		Protocol:     "tcp",
	}
}

func TestCorrelatorFirstObservationNotCorrelated(t *testing.T) {
	cfg := CorrelatorConfig{WindowDuration: 10 * time.Second, MinOccurrences: 2}
	c := NewCorrelator(cfg)
	p := makeCorrelatorPort(8080)
	if c.Observe(p) {
		t.Error("expected first observation to not be correlated")
	}
}

func TestCorrelatorReachesMinOccurrences(t *testing.T) {
	cfg := CorrelatorConfig{WindowDuration: 10 * time.Second, MinOccurrences: 3}
	c := NewCorrelator(cfg)
	p := makeCorrelatorPort(9090)
	c.Observe(p)
	c.Observe(p)
	if !c.Observe(p) {
		t.Error("expected third observation to be correlated")
	}
}

func TestCorrelatorWindowExpiry(t *testing.T) {
	now := time.Now()
	cfg := CorrelatorConfig{WindowDuration: 5 * time.Second, MinOccurrences: 2}
	c := NewCorrelator(cfg)
	c.now = func() time.Time { return now }
	p := makeCorrelatorPort(7070)

	// First observation at t=0.
	c.Observe(p)

	// Advance time past the window so the first observation expires.
	c.now = func() time.Time { return now.Add(6 * time.Second) }

	// Second observation should not correlate because first is expired.
	if c.Observe(p) {
		t.Error("expected observation after window expiry to not correlate")
	}
}

func TestCorrelatorIndependentPorts(t *testing.T) {
	cfg := CorrelatorConfig{WindowDuration: 10 * time.Second, MinOccurrences: 2}
	c := NewCorrelator(cfg)
	p1 := makeCorrelatorPort(1111)
	p2 := makeCorrelatorPort(2222)

	c.Observe(p1)
	c.Observe(p1)
	if c.Observe(p2) {
		t.Error("p2 should not be correlated after only one observation")
	}
}

func TestCorrelatorReset(t *testing.T) {
	cfg := CorrelatorConfig{WindowDuration: 10 * time.Second, MinOccurrences: 2}
	c := NewCorrelator(cfg)
	p := makeCorrelatorPort(3333)
	c.Observe(p)
	c.Reset(p)
	if c.Observe(p) {
		t.Error("expected observation after reset to not correlate")
	}
}

func TestCorrelatorPurge(t *testing.T) {
	now := time.Now()
	cfg := CorrelatorConfig{WindowDuration: 5 * time.Second, MinOccurrences: 2}
	c := NewCorrelator(cfg)
	c.now = func() time.Time { return now }
	p := makeCorrelatorPort(4444)
	c.Observe(p)

	c.now = func() time.Time { return now.Add(10 * time.Second) }
	c.Purge()

	if len(c.events) != 0 {
		t.Errorf("expected events map to be empty after purge, got %d entries", len(c.events))
	}
}
