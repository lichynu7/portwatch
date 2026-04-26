package ports

import (
	"testing"
)

func makeRawPort(port uint16, proto string) Port {
	return Port{
		Port:     port,
		Protocol: proto,
		Address:  "0.0.0.0",
		Inode:    0,
	}
}

func TestEnricherNoOptions(t *testing.T) {
	e := NewEnricher()
	ports := []Port{makeRawPort(80, "tcp"), makeRawPort(443, "tcp")}
	result := e.Enrich(ports)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	for _, r := range result {
		if r.Service != "" {
			t.Errorf("expected empty service without resolution, got %q", r.Service)
		}
		if r.Process != nil {
			t.Errorf("expected nil process without lookup")
		}
	}
}

func TestEnricherServiceResolution(t *testing.T) {
	e := NewEnricher(WithServiceResolution())
	ports := []Port{
		makeRawPort(22, "tcp"),
		makeRawPort(80, "tcp"),
		makeRawPort(443, "tcp"),
		makeRawPort(9999, "tcp"),
	}
	result := e.Enrich(ports)

	expected := map[uint16]string{
		22:   "ssh",
		80:   "http",
		443:  "https",
		9999: "",
	}
	for _, r := range result {
		want, ok := expected[r.Port]
		if !ok {
			continue
		}
		if r.Service != want {
			t.Errorf("port %d: expected service %q, got %q", r.Port, want, r.Service)
		}
	}
}

func TestEnricherPreservesFields(t *testing.T) {
	e := NewEnricher()
	p := makeRawPort(8080, "tcp")
	p.Address = "127.0.0.1"
	result := e.Enrich([]Port{p})
	if len(result) != 1 {
		t.Fatal("expected 1 result")
	}
	if result[0].Port != 8080 {
		t.Errorf("port mismatch: got %d", result[0].Port)
	}
	if result[0].Address != "127.0.0.1" {
		t.Errorf("address mismatch: got %s", result[0].Address)
	}
	if result[0].Protocol != "tcp" {
		t.Errorf("protocol mismatch: got %s", result[0].Protocol)
	}
}

func TestEnricherEmptyInput(t *testing.T) {
	e := NewEnricher(WithServiceResolution())
	result := e.Enrich([]Port{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestKnownServicesMap(t *testing.T) {
	cases := []struct {
		port    uint16
		service string
	}{
		{22, "ssh"},
		{3306, "mysql"},
		{5432, "postgres"},
		{6379, "redis"},
	}
	for _, tc := range cases {
		got := knownServices[tc.port]
		if got != tc.service {
			t.Errorf("port %d: expected %q, got %q", tc.port, tc.service, got)
		}
	}
}
