package ports

import (
	"testing"
	"time"
)

func TestThrottleAllowsFirstAlert(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestThrottleSuppressesWithinWindow(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Allow("tcp:8080") // first — allowed
	if th.Allow("tcp:8080") {
		t.Fatal("expected second alert within window to be suppressed")
	}
}

func TestThrottleAllowsAfterWindowExpires(t *testing.T) {
	cfg := ThrottleConfig{Window: 10 * time.Millisecond, MaxBurst: 1}
	th := NewThrottle(cfg)

	var fakeNow time.Time
	th.now = func() time.Time { return fakeNow }

	th.Allow("tcp:9090")
	fakeNow = fakeNow.Add(20 * time.Millisecond)

	if !th.Allow("tcp:9090") {
		t.Fatal("expected alert after window expiry to be allowed")
	}
}

func TestThrottleIndependentKeys(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Allow("tcp:8080")

	if !th.Allow("tcp:9090") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestThrottleBurstGreaterThanOne(t *testing.T) {
	cfg := ThrottleConfig{Window: 5 * time.Minute, MaxBurst: 3}
	th := NewThrottle(cfg)

	for i := 0; i < 3; i++ {
		if !th.Allow("tcp:443") {
			t.Fatalf("expected alert %d within burst to be allowed", i+1)
		}
	}
	if th.Allow("tcp:443") {
		t.Fatal("expected alert beyond burst to be suppressed")
	}
}

func TestThrottleResetAllowsNext(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Allow("tcp:8080")
	th.Reset("tcp:8080")

	if !th.Allow("tcp:8080") {
		t.Fatal("expected alert after Reset to be allowed")
	}
}

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.Window <= 0 {
		t.Error("expected positive window duration")
	}
	if cfg.MaxBurst < 1 {
		t.Error("expected MaxBurst >= 1")
	}
}
