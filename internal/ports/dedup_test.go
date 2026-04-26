package ports

import (
	"testing"
	"time"
)

func TestDedupFirstAlertAllowed(t *testing.T) {
	d := NewDeduplicator(DefaultDedupConfig())
	if d.IsDuplicate("tcp:8080") {
		t.Fatal("expected first occurrence to not be a duplicate")
	}
}

func TestDedupSecondAlertSuppressed(t *testing.T) {
	d := NewDeduplicator(DefaultDedupConfig())
	d.IsDuplicate("tcp:8080")
	if !d.IsDuplicate("tcp:8080") {
		t.Fatal("expected second occurrence to be a duplicate")
	}
}

func TestDedupAllowsAfterWindowExpires(t *testing.T) {
	cfg := DedupConfig{WindowSize: 100 * time.Millisecond}
	d := NewDeduplicator(cfg)

	fake := time.Now()
	d.now = func() time.Time { return fake }

	d.IsDuplicate("tcp:9090")

	// advance time beyond window
	fake = fake.Add(200 * time.Millisecond)
	if d.IsDuplicate("tcp:9090") {
		t.Fatal("expected alert after window expiry to not be duplicate")
	}
}

func TestDedupIndependentKeys(t *testing.T) {
	d := NewDeduplicator(DefaultDedupConfig())
	d.IsDuplicate("tcp:8080")
	if d.IsDuplicate("tcp:9090") {
		t.Fatal("expected different key to not be a duplicate")
	}
}

func TestDedupPurgeRemovesExpired(t *testing.T) {
	cfg := DedupConfig{WindowSize: 50 * time.Millisecond}
	d := NewDeduplicator(cfg)

	fake := time.Now()
	d.now = func() time.Time { return fake }

	d.IsDuplicate("tcp:1234")

	fake = fake.Add(100 * time.Millisecond)
	d.Purge()

	// After purge, key should be fresh again
	if d.IsDuplicate("tcp:1234") {
		t.Fatal("expected purged key to not be a duplicate")
	}
}

func TestDedupReset(t *testing.T) {
	d := NewDeduplicator(DefaultDedupConfig())
	d.IsDuplicate("tcp:8080")
	d.Reset()
	if d.IsDuplicate("tcp:8080") {
		t.Fatal("expected key to be cleared after reset")
	}
}

func TestDefaultDedupConfig(t *testing.T) {
	cfg := DefaultDedupConfig()
	if cfg.WindowSize <= 0 {
		t.Fatalf("expected positive WindowSize, got %v", cfg.WindowSize)
	}
}
