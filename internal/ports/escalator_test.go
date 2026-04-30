package ports

import (
	"testing"
	"time"
)

func TestEscalatorBelowThreshold(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 3
	e := NewEscalator(cfg)

	if e.Observe("tcp:8080") {
		t.Fatal("expected no escalation on first observation")
	}
	if e.Observe("tcp:8080") {
		t.Fatal("expected no escalation on second observation")
	}
}

func TestEscalatorReachesThreshold(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 3
	e := NewEscalator(cfg)

	e.Observe("tcp:8080")
	e.Observe("tcp:8080")
	if !e.Observe("tcp:8080") {
		t.Fatal("expected escalation on third observation")
	}
}

func TestEscalatorWindowExpiry(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 2
	cfg.Window = 100 * time.Millisecond
	e := NewEscalator(cfg)

	fixed := time.Now()
	e.now = func() time.Time { return fixed }
	e.Observe("tcp:9090")

	// Advance past the window.
	e.now = func() time.Time { return fixed.Add(200 * time.Millisecond) }
	if e.Observe("tcp:9090") {
		t.Fatal("expected window reset, no escalation")
	}
}

func TestEscalatorIndependentKeys(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 2
	e := NewEscalator(cfg)

	e.Observe("tcp:1111")
	if e.Observe("tcp:2222") {
		t.Fatal("second key should not escalate independently")
	}
}

func TestEscalatorDisabled(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Enabled = false
	cfg.Threshold = 1
	e := NewEscalator(cfg)

	if e.Observe("tcp:8080") {
		t.Fatal("disabled escalator should never escalate")
	}
}

func TestEscalatorReset(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 2
	e := NewEscalator(cfg)

	e.Observe("tcp:7070")
	e.Reset("tcp:7070")
	if e.Observe("tcp:7070") {
		t.Fatal("expected no escalation after reset")
	}
}

func TestEscalatorPurge(t *testing.T) {
	cfg := DefaultEscalatorConfig()
	cfg.Threshold = 5
	cfg.Window = 50 * time.Millisecond
	e := NewEscalator(cfg)

	fixed := time.Now()
	e.now = func() time.Time { return fixed }
	e.Observe("tcp:3000")

	e.now = func() time.Time { return fixed.Add(100 * time.Millisecond) }
	e.Purge()

	e.mu.Lock()
	_, exists := e.entries["tcp:3000"]
	e.mu.Unlock()
	if exists {
		t.Fatal("expected expired entry to be purged")
	}
}
