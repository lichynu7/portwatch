package ports

import (
	"testing"
	"time"
)

func TestDecayDisabledReturnsRaw(t *testing.T) {
	cfg := DecayConfig{Enabled: false}
	d := NewScoreDecayer(cfg)
	got := d.Update("tcp:8080", 42.0)
	if got != 42.0 {
		t.Fatalf("expected 42.0, got %f", got)
	}
}

func TestDecayFirstUpdateNoDecay(t *testing.T) {
	cfg := DefaultDecayConfig()
	d := NewScoreDecayer(cfg)
	got := d.Update("tcp:8080", 10.0)
	if got != 10.0 {
		t.Fatalf("expected 10.0, got %f", got)
	}
}

func TestDecayAccumulatesWithinHalfLife(t *testing.T) {
	now := time.Now()
	cfg := DefaultDecayConfig()
	d := NewScoreDecayer(cfg)
	d.now = func() time.Time { return now }

	d.Update("tcp:9000", 8.0)
	// advance by less than one half-life (5 min < 10 min)
	d.now = func() time.Time { return now.Add(5 * time.Minute) }
	got := d.Update("tcp:9000", 4.0)
	// after 5 min (0.5 half-lives): 8 / 2^0.5 ≈ 5.657, plus 4 ≈ 9.657
	if got < 9.0 || got > 10.5 {
		t.Fatalf("expected ~9.66, got %f", got)
	}
}

func TestDecayAfterFullHalfLife(t *testing.T) {
	now := time.Now()
	cfg := DefaultDecayConfig()
	d := NewScoreDecayer(cfg)
	d.now = func() time.Time { return now }

	d.Update("tcp:443", 20.0)
	d.now = func() time.Time { return now.Add(10 * time.Minute) }
	got := d.Update("tcp:443", 0.0)
	// after one half-life: 20 / 2 = 10
	if got < 9.9 || got > 10.1 {
		t.Fatalf("expected ~10.0, got %f", got)
	}
}

func TestDecayGetReturnsDecayedValue(t *testing.T) {
	now := time.Now()
	cfg := DefaultDecayConfig()
	d := NewScoreDecayer(cfg)
	d.now = func() time.Time { return now }

	d.Update("tcp:22", 16.0)
	d.now = func() time.Time { return now.Add(20 * time.Minute) } // two half-lives
	got := d.Get("tcp:22")
	// 16 / 2^2 = 4
	if got < 3.9 || got > 4.1 {
		t.Fatalf("expected ~4.0, got %f", got)
	}
}

func TestDecayGetMissingKeyReturnsZero(t *testing.T) {
	d := NewScoreDecayer(DefaultDecayConfig())
	if got := d.Get("unknown"); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestDecayPurgeRemovesLowEntries(t *testing.T) {
	now := time.Now()
	cfg := DefaultDecayConfig()
	cfg.MinScore = 5.0
	d := NewScoreDecayer(cfg)
	d.now = func() time.Time { return now }

	d.Update("tcp:80", 4.0) // starts below MinScore threshold after decay
	d.Update("tcp:443", 50.0)

	// advance far enough that tcp:80 decays below MinScore
	d.now = func() time.Time { return now.Add(30 * time.Minute) }
	d.Purge()

	if got := d.Get("tcp:80"); got != 0 {
		t.Fatalf("expected tcp:80 to be purged, got %f", got)
	}
	if got := d.Get("tcp:443"); got == 0 {
		t.Fatalf("expected tcp:443 to survive purge")
	}
}

func TestDecayIndependentKeys(t *testing.T) {
	now := time.Now()
	cfg := DefaultDecayConfig()
	d := NewScoreDecayer(cfg)
	d.now = func() time.Time { return now }

	d.Update("tcp:8080", 10.0)
	d.Update("udp:53", 5.0)

	if a, b := d.Get("tcp:8080"), d.Get("udp:53"); a == b {
		t.Fatalf("expected independent scores, both are %f", a)
	}
}
