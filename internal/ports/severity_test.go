package ports

import (
	"testing"
)

func TestSeverityString(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityInfo, "info"},
		{SeverityWarning, "warning"},
		{SeverityCritical, "critical"},
		{Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q; want %q", tc.s, got, tc.want)
		}
	}
}

func TestClassifyPrivilegedPort(t *testing.T) {
	cfg := DefaultSeverityConfig()
	p := Port{Port: 22, Process: "sshd"}
	if got := Classify(p, cfg); got != SeverityCritical {
		t.Errorf("port 22 should be Critical, got %s", got)
	}
}

func TestClassifyEphemeralPort(t *testing.T) {
	cfg := DefaultSeverityConfig()
	p := Port{Port: 45000, Process: "curl"}
	if got := Classify(p, cfg); got != SeverityInfo {
		t.Errorf("port 45000 should be Info, got %s", got)
	}
}

func TestClassifyWarningPort(t *testing.T) {
	cfg := DefaultSeverityConfig()
	p := Port{Port: 8080, Process: "app"}
	if got := Classify(p, cfg); got != SeverityWarning {
		t.Errorf("port 8080 should be Warning, got %s", got)
	}
}

func TestClassifyCriticalProcess(t *testing.T) {
	cfg := DefaultSeverityConfig()
	// High port but process name triggers critical.
	p := Port{Port: 9999, Process: "ncat"}
	if got := Classify(p, cfg); got != SeverityCritical {
		t.Errorf("process 'ncat' on port 9999 should be Critical, got %s", got)
	}
}

func TestClassifyCriticalProcessCaseInsensitive(t *testing.T) {
	cfg := DefaultSeverityConfig()
	p := Port{Port: 12345, Process: "SOCAT"}
	if got := Classify(p, cfg); got != SeverityCritical {
		t.Errorf("process 'SOCAT' should be Critical regardless of case, got %s", got)
	}
}

func TestClassifyZeroPort(t *testing.T) {
	cfg := DefaultSeverityConfig()
	// Port 0 is not in privileged range (> 0 required), falls to Warning.
	p := Port{Port: 0, Process: "unknown"}
	if got := Classify(p, cfg); got != SeverityWarning {
		t.Errorf("port 0 should be Warning, got %s", got)
	}
}
