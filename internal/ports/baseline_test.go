package ports

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaselineEmpty(t *testing.T) {
	b := NewBaseline("/tmp/nonexistent.json")
	if len(b.Entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(b.Entries))
	}
}

func TestBaselineContains(t *testing.T) {
	b := NewBaseline("/tmp/nonexistent.json")
	if b.Contains("tcp", 80) {
		t.Fatal("expected false for empty baseline")
	}
}

func TestBaselineAddAndContains(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	b := NewBaseline(path)

	entry := BaselineEntry{Port: 8080, Proto: "tcp", Process: "nginx", Reason: "known web server"}
	if err := b.Add(entry); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !b.Contains("tcp", 8080) {
		t.Fatal("expected baseline to contain tcp:8080")
	}
	if b.Contains("udp", 8080) {
		t.Fatal("proto mismatch should not match")
	}
}

func TestBaselinePersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := NewBaseline(path)
	_ = b.Add(BaselineEntry{Port: 443, Proto: "tcp", Process: "nginx", Reason: "https"})

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if !loaded.Contains("tcp", 443) {
		t.Fatal("expected loaded baseline to contain tcp:443")
	}
}

func TestLoadBaselineMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")
	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(b.Entries) != 0 {
		t.Fatal("expected empty baseline for missing file")
	}
}

func TestBaselineRemove(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	b := NewBaseline(path)
	_ = b.Add(BaselineEntry{Port: 22, Proto: "tcp", Process: "sshd", Reason: "ssh"})

	if err := b.Remove("tcp", 22); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if b.Contains("tcp", 22) {
		t.Fatal("expected port removed from baseline")
	}

	// Verify persisted
	loaded, _ := LoadBaseline(path)
	if loaded.Contains("tcp", 22) {
		t.Fatal("expected removed port absent after reload")
	}
}

func TestLoadBaselineCorrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o640)
	_, err := LoadBaseline(path)
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}
