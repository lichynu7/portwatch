package ports

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFakeTCP(t *testing.T, dir, content string) {
	t.Helper()
	net4 := filepath.Join(dir, "net")
	_ = os.MkdirAll(net4, 0o755)
	err := os.WriteFile(filepath.Join(net4, "tcp"), []byte(content), 0o644)
	if err != nil {
		t.Fatalf("writeFakeTCP: %v", err)
	}
}

const tcpHeader = "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"

func TestNewWatcherInvalidInterval(t *testing.T) {
	_, err := NewWatcher(WatchConfig{Interval: 0})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewWatcherValid(t *testing.T) {
	w, err := NewWatcher(WatchConfig{Interval: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestWatcherRunCancels(t *testing.T) {
	dir := t.TempDir()
	writeFakeTCP(t, dir, tcpHeader)

	w, err := NewWatcher(WatchConfig{
		Interval: 10 * time.Millisecond,
		ProcFS:   dir,
	})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = w.Run(ctx)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestWatcherDetectsNewPort(t *testing.T) {
	dir := t.TempDir()
	writeFakeTCP(t, dir, tcpHeader)

	w, err := NewWatcher(WatchConfig{
		Interval: 20 * time.Millisecond,
		ProcFS:   dir,
	})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Run(ctx) //nolint:errcheck

	// After initial scan, inject a new port.
	time.Sleep(30 * time.Millisecond)
	writeFakeTCP(t, dir, tcpHeader+"   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1\n")

	select {
	case ev := <-w.Events():
		if ev.EventType != "new" {
			t.Fatalf("expected 'new', got %q", ev.EventType)
		}
		if ev.Port.Port != 8080 {
			t.Fatalf("expected port 8080, got %d", ev.Port.Port)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for port event")
	}
}

func TestIndexPorts(t *testing.T) {
	ports := []Port{
		{Proto: "tcp", Port: 80},
		{Proto: "tcp", Port: 443},
	}
	m := indexPorts(ports)
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if _, ok := m["tcp:80"]; !ok {
		t.Error("missing tcp:80")
	}
}
