package ports

import (
	"testing"
)

func makePortSlice(specs []struct {
	proto string
	port  uint16
	state string
}) []Port {
	out := make([]Port, 0, len(specs))
	for _, s := range specs {
		out = append(out, Port{Protocol: s.proto, Port: s.port, State: s.state})
	}
	return out
}

func TestFingerprintEmpty(t *testing.T) {
	f := NewFingerprint(nil)
	if f.Count != 0 {
		t.Fatalf("expected count 0, got %d", f.Count)
	}
	if f.Hash == "" {
		t.Fatal("expected non-empty hash for empty set")
	}
}

func TestFingerprintDeterministic(t *testing.T) {
	ports := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{
		{"tcp", 80, "LISTEN"},
		{"tcp", 443, "LISTEN"},
	})
	f1 := NewFingerprint(ports)
	f2 := NewFingerprint(ports)
	if !f1.Equal(f2) {
		t.Fatalf("expected identical fingerprints, got %s vs %s", f1.Hash, f2.Hash)
	}
}

func TestFingerprintOrderIndependent(t *testing.T) {
	spec := []struct {
		proto string
		port  uint16
		state string
	}{
		{"tcp", 8080, "LISTEN"},
		{"tcp", 22, "LISTEN"},
	}
	a := makePortSlice(spec)
	b := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{
		{"tcp", 22, "LISTEN"},
		{"tcp", 8080, "LISTEN"},
	})
	if !NewFingerprint(a).Equal(NewFingerprint(b)) {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestFingerprintDiffersOnChange(t *testing.T) {
	a := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{{"tcp", 80, "LISTEN"}})
	b := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{{"tcp", 9000, "LISTEN"}})
	if NewFingerprint(a).Equal(NewFingerprint(b)) {
		t.Fatal("different port sets should produce different fingerprints")
	}
}

func TestFingerprintString(t *testing.T) {
	ports := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{{"tcp", 443, "LISTEN"}})
	s := NewFingerprint(ports).String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
}

func TestFingerprintCount(t *testing.T) {
	ports := makePortSlice([]struct {
		proto string
		port  uint16
		state string
	}{
		{"tcp", 80, "LISTEN"},
		{"tcp", 443, "LISTEN"},
		{"udp", 53, "LISTEN"},
	})
	f := NewFingerprint(ports)
	if f.Count != 3 {
		t.Fatalf("expected count 3, got %d", f.Count)
	}
}
