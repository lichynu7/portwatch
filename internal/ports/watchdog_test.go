package ports

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatchdogDisabledReturnsNil(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	cfg.Enabled = false
	wd := NewWatchdog(cfg)
	if wd != nil {
		t.Fatal("expected nil watchdog when disabled")
	}
}

func TestWatchdogDefaultConfig(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled=true by default")
	}
	if cfg.StallTimeout <= 0 {
		t.Error("expected positive StallTimeout")
	}
}

func TestWatchdogTickPreventsStall(t *testing.T) {
	var fired atomic.Bool
	cfg := WatchdogConfig{
		Enabled:      true,
		StallTimeout: 200 * time.Millisecond,
		OnStall: func(_ time.Duration) {
			fired.Store(true)
		},
	}
	wd := NewWatchdog(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Run(ctx)

	// Tick frequently — stall should not fire.
	for i := 0; i < 5; i++ {
		time.Sleep(30 * time.Millisecond)
		wd.Tick()
	}
	if fired.Load() {
		t.Error("watchdog fired despite regular ticks")
	}
}

func TestWatchdogFiresOnStall(t *testing.T) {
	var fired atomic.Bool
	cfg := WatchdogConfig{
		Enabled:      true,
		StallTimeout: 80 * time.Millisecond,
		OnStall: func(_ time.Duration) {
			fired.Store(true)
		},
	}
	wd := NewWatchdog(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Run(ctx)

	// Do NOT tick — watchdog should fire after StallTimeout.
	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if fired.Load() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("watchdog did not fire after stall timeout")
}

func TestWatchdogRunCancels(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	wd := NewWatchdog(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		wd.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}
