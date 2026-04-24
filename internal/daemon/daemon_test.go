package daemon

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.Interval = 50 * time.Millisecond
	cfg.SnapshotPath = filepath.Join(t.TempDir(), "snap.json")
	return cfg
}

func TestDaemonRunCancels(t *testing.T) {
	cfg := testConfig(t)
	disp := alert.NewDispatcher()
	d := New(cfg, disp)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got %v", err)
	}
}

func TestDaemonCreatesSnapshot(t *testing.T) {
	cfg := testConfig(t)
	disp := alert.NewDispatcher()
	d := New(cfg, disp)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = d.Run(ctx)

	if _, err := os.Stat(cfg.SnapshotPath); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to be created")
	}
}

func TestDaemonAlertOnNewPort(t *testing.T) {
	cfg := testConfig(t)

	received := make([]alert.Alert, 0)
	handler := alert.HandlerFunc(func(_ context.Context, a alert.Alert) error {
		received = append(received, a)
		return nil
	})

	disp := alert.NewDispatcher()
	disp.Register(handler)

	d := New(cfg, disp)

	// Run one tick to create baseline snapshot.
	_ = d.tick()

	// Verify no alerts on first tick (baseline).
	if len(received) != 0 {
		t.Fatalf("expected 0 alerts on baseline tick, got %d", len(received))
	}
}
