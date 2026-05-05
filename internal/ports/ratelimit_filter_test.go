package ports

import (
	"testing"
	"time"
)

func TestRateLimitFilterAllowsUpToBurst(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 10 * time.Second, MaxHits: 3}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	for i := 0; i < 3; i++ {
		if !f.Allow("tcp:8080", now) {
			t.Fatalf("hit %d: expected allow", i+1)
		}
	}
}

func TestRateLimitFilterSuppressesAboveBurst(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 10 * time.Second, MaxHits: 2}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	f.Allow("tcp:9000", now)
	f.Allow("tcp:9000", now)
	if f.Allow("tcp:9000", now) {
		t.Fatal("expected suppression on third hit")
	}
}

func TestRateLimitFilterResetsAfterWindow(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 5 * time.Millisecond, MaxHits: 1}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	f.Allow("tcp:7070", now)
	if f.Allow("tcp:7070", now) {
		t.Fatal("expected suppression within window")
	}

	later := now.Add(10 * time.Millisecond)
	if !f.Allow("tcp:7070", later) {
		t.Fatal("expected allow after window expired")
	}
}

func TestRateLimitFilterIndependentKeys(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 10 * time.Second, MaxHits: 1}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	if !f.Allow("tcp:1111", now) {
		t.Fatal("key1 first hit should be allowed")
	}
	if !f.Allow("tcp:2222", now) {
		t.Fatal("key2 first hit should be allowed")
	}
	if f.Allow("tcp:1111", now) {
		t.Fatal("key1 second hit should be suppressed")
	}
}

func TestRateLimitFilterReset(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 10 * time.Second, MaxHits: 1}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	f.Allow("tcp:3333", now)
	f.Reset("tcp:3333")
	if !f.Allow("tcp:3333", now) {
		t.Fatal("expected allow after explicit reset")
	}
}

func TestApplyRateLimitFilter(t *testing.T) {
	cfg := RateLimitFilterConfig{Window: 10 * time.Second, MaxHits: 1}
	f := NewAlertRateLimitFilter(cfg)
	now := time.Now()

	ports := []Port{
		{Port: 8080, Protocol: "tcp"},
		{Port: 8080, Protocol: "tcp"},
		{Port: 9090, Protocol: "tcp"},
	}

	out := ApplyRateLimitFilter(ports, f, now)
	if len(out) != 2 {
		t.Fatalf("expected 2 ports after filter, got %d", len(out))
	}
}
