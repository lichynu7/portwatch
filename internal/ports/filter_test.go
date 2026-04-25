package ports

import (
	"testing"
)

func makePort(port uint16, process string) Port {
	return Port{Port: port, Process: process}
}

func TestNewFilterConfig(t *testing.T) {
	fc := NewFilterConfig([]uint16{22, 80}, []string{"nginx", "sshd"})
	if !fc.AllowedPorts[22] || !fc.AllowedPorts[80] {
		t.Error("expected ports 22 and 80 to be in AllowedPorts")
	}
	if len(fc.AllowedProcesses) != 2 {
		t.Errorf("expected 2 allowed processes, got %d", len(fc.AllowedProcesses))
	}
}

func TestIsSafeByPort(t *testing.T) {
	fc := NewFilterConfig([]uint16{443}, nil)
	if !fc.IsSafe(makePort(443, "unknown")) {
		t.Error("port 443 should be safe")
	}
	if fc.IsSafe(makePort(9999, "unknown")) {
		t.Error("port 9999 should not be safe")
	}
}

func TestIsSafeByProcess(t *testing.T) {
	fc := NewFilterConfig(nil, []string{"nginx", "sshd"})
	if !fc.IsSafe(makePort(12345, "nginx: worker")) {
		t.Error("nginx process should be safe")
	}
	if fc.IsSafe(makePort(12345, "malware")) {
		t.Error("malware process should not be safe")
	}
}

func TestIsSafeNilConfig(t *testing.T) {
	var fc *FilterConfig
	if fc.IsSafe(makePort(80, "nginx")) {
		t.Error("nil FilterConfig should never mark a port as safe")
	}
}

func TestExcludeSafe(t *testing.T) {
	ports := []Port{
		makePort(22, "sshd"),
		makePort(8080, "unknown-app"),
		makePort(443, "nginx"),
	}
	fc := NewFilterConfig([]uint16{22, 443}, nil)
	result := ExcludeSafe(ports, fc)
	if len(result) != 1 {
		t.Fatalf("expected 1 unexpected port, got %d", len(result))
	}
	if result[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", result[0].Port)
	}
}

func TestExcludeSafeNilConfig(t *testing.T) {
	ports := []Port{
		makePort(22, "sshd"),
		makePort(8080, "app"),
	}
	result := ExcludeSafe(ports, nil)
	if len(result) != len(ports) {
		t.Errorf("nil config should return all ports unchanged, got %d", len(result))
	}
}
