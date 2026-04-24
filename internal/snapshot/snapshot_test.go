package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makePorts(specs ...string) []ports.Port {
	var ps []ports.Port
	for i, s := range specs {
		ps = append(ps, ports.Port{Protocol: s, LocalPort: uint16(8000 + i)})
	}
	return ps
}

func TestNewSnapshot(t *testing.T) {
	ps := makePorts("tcp", "tcp")
	s := New(ps)
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(s.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(s.Ports))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := New(makePorts("tcp", "udp"))
	orig.Timestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Timestamp.Equal(orig.Timestamp) {
		t.Errorf("timestamp mismatch: got %v, want %v", loaded.Timestamp, orig.Timestamp)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch: got %d, want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoadMissingFile(t *testing.T) {
	s, err := Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if s == nil || len(s.Ports) != 0 {
		t.Error("expected empty snapshot for missing file")
	}
}

func TestDiff(t *testing.T) {
	prev := New([]ports.Port{
		{Protocol: "tcp", LocalPort: 80},
		{Protocol: "tcp", LocalPort: 443},
	})
	curr := New([]ports.Port{
		{Protocol: "tcp", LocalPort: 443},
		{Protocol: "tcp", LocalPort: 8080},
	})

	added, removed := Diff(prev, curr)

	if len(added) != 1 || added[0].LocalPort != 8080 {
		t.Errorf("unexpected added ports: %+v", added)
	}
	if len(removed) != 1 || removed[0].LocalPort != 80 {
		t.Errorf("unexpected removed ports: %+v", removed)
	}
}

func TestDiffNoChange(t *testing.T) {
	ps := makePorts("tcp", "udp")
	added, removed := Diff(New(ps), New(ps))
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestSaveInvalidPath(t *testing.T) {
	s := New(makePorts("tcp"))
	err := s.Save(filepath.Join(os.DevNull, "nope", "snap.json"))
	if err == nil {
		t.Error("expected error saving to invalid path")
	}
}
