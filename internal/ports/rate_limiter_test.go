package ports

import (
	"testing"
	"time"
)

func TestRateLimiterAllowsFirstAlert(t *testing.T) {
	rl := NewRateLimiter(5 * time.Minute)
	if !rl.Allow(8080) {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestRateLimiterSuppressesDuplicate(t *testing.T) {
	rl := NewRateLimiter(5 * time.Minute)
	rl.Allow(8080) // first — records timestamp
	if rl.Allow(8080) {
		t.Fatal("expected second alert within cooldown to be suppressed")
	}
}

func TestRateLimiterAllowsAfterCooldown(t *testing.T) {
	rl := NewRateLimiter(10 * time.Millisecond)
	rl.Allow(9000)
	time.Sleep(20 * time.Millisecond)
	if !rl.Allow(9000) {
		t.Fatal("expected alert to be allowed after cooldown expired")
	}
}

func TestRateLimiterIndependentPorts(t *testing.T) {
	rl := NewRateLimiter(5 * time.Minute)
	rl.Allow(80)
	if !rl.Allow(443) {
		t.Fatal("suppressing port 80 should not affect port 443")
	}
}

func TestRateLimiterReset(t *testing.T) {
	rl := NewRateLimiter(5 * time.Minute)
	rl.Allow(3000)
	rl.Reset(3000)
	if !rl.Allow(3000) {
		t.Fatal("expected alert to be allowed after explicit reset")
	}
}

func TestRateLimiterPurge(t *testing.T) {
	rl := NewRateLimiter(10 * time.Millisecond)
	rl.Allow(1234)
	rl.Allow(5678)
	time.Sleep(20 * time.Millisecond)
	rl.Purge()

	rl.mu.Lock()
	remaining := len(rl.lastSeen)
	rl.mu.Unlock()

	if remaining != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", remaining)
	}
}

func TestRateLimiterPurgeKeepsRecent(t *testing.T) {
	rl := NewRateLimiter(5 * time.Minute)
	rl.Allow(7777)
	rl.Purge() // cooldown has not elapsed

	rl.mu.Lock()
	remaining := len(rl.lastSeen)
	rl.mu.Unlock()

	if remaining != 1 {
		t.Fatalf("expected recent entry to survive purge, got %d entries", remaining)
	}
}
