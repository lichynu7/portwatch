package ports

import (
	"path/filepath"
	"testing"
)

func makeTestPort(proto string, port uint16) Port {
	return Port{Proto: proto, Port: port, LocalAddr: "0.0.0.0"}
}

func TestApplyBaselineNilBaseline(t *testing.T) {
	ports := []Port{
		makeTestPort("tcp", 80),
		makeTestPort("tcp", 443),
	}
	result := ApplyBaseline(ports, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports with nil baseline, got %d", len(result))
	}
}

func TestApplyBaselineFiltersKnown(t *testing.T) {
	dir := t.TempDir()
	b := NewBaseline(filepath.Join(dir, "b.json"))
	_ = b.Add(BaselineEntry{Port: 22, Proto: "tcp"})
	_ = b.Add(BaselineEntry{Port: 80, Proto: "tcp"})

	ports := []Port{
		makeTestPort("tcp", 22),
		makeTestPort("tcp", 80),
		makeTestPort("tcp", 9999),
	}
	result := ApplyBaseline(ports, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 unexpected port, got %d", len(result))
	}
	if result[0].Port != 9999 {
		t.Fatalf("expected port 9999, got %d", result[0].Port)
	}
}

func TestApplyBaselineEmptyBaseline(t *testing.T) {
	dir := t.TempDir()
	b := NewBaseline(filepath.Join(dir, "b.json"))

	ports := []Port{
		makeTestPort("tcp", 8080),
	}
	result := ApplyBaseline(ports, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
}

func TestExcludeBaselineFilterFunc(t *testing.T) {
	dir := t.TempDir()
	b := NewBaseline(filepath.Join(dir, "b.json"))
	_ = b.Add(BaselineEntry{Port: 3306, Proto: "tcp"})

	f := ExcludeBaseline(b)

	if !f(makeTestPort("tcp", 3306)) {
		t.Fatal("expected filter to return true (is safe/known) for baselined port")
	}
	if f(makeTestPort("tcp", 5432)) {
		t.Fatal("expected filter to return false for unknown port")
	}
}

func TestApplyBaselineProtoSensitive(t *testing.T) {
	dir := t.TempDir()
	b := NewBaseline(filepath.Join(dir, "b.json"))
	_ = b.Add(BaselineEntry{Port: 53, Proto: "tcp"})

	ports := []Port{
		makeTestPort("tcp", 53),
		makeTestPort("udp", 53),
	}
	result := ApplyBaseline(ports, b)
	// udp:53 should NOT be filtered — only tcp:53 is baselined
	if len(result) != 1 {
		t.Fatalf("expected 1 port (udp:53), got %d", len(result))
	}
	if result[0].Proto != "udp" {
		t.Fatalf("expected udp proto, got %s", result[0].Proto)
	}
}
