package ports

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func defaultLabelerCfg() config.LabelerConfig {
	return config.LabelerConfig{
		TrustedProcesses: []string{"sshd", "nginx"},
	}
}

func TestNewLabelerValid(t *testing.T) {
	_, err := NewLabeler(defaultLabelerCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLabelProtocol(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Protocol: "TCP", Port: 80}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelProtocol]; got != "tcp" {
		t.Errorf("expected tcp, got %q", got)
	}
}

func TestLabelPortRangePrivileged(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 22, Protocol: "tcp"}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelPortRange]; got != "privileged" {
		t.Errorf("expected privileged, got %q", got)
	}
}

func TestLabelPortRangeRegistered(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 8080, Protocol: "tcp"}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelPortRange]; got != "registered" {
		t.Errorf("expected registered, got %q", got)
	}
}

func TestLabelPortRangeEphemeral(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 55000, Protocol: "tcp"}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelPortRange]; got != "ephemeral" {
		t.Errorf("expected ephemeral, got %q", got)
	}
}

func TestLabelTrustedProcess(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 22, Protocol: "tcp", Process: "sshd"}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelTrusted]; got != "true" {
		t.Errorf("expected trusted=true, got %q", got)
	}
}

func TestLabelUntrustedProcess(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 4444, Protocol: "tcp", Process: "nc"}}
	out := l.Label(ports)
	if got := out[0].Labels[LabelTrusted]; got != "false" {
		t.Errorf("expected trusted=false, got %q", got)
	}
}

func TestLabelInitialisesNilMap(t *testing.T) {
	l, _ := NewLabeler(defaultLabelerCfg())
	ports := []Port{{Port: 80, Protocol: "tcp"}}
	if ports[0].Labels != nil {
		t.Fatal("expected nil labels before Label call")
	}
	out := l.Label(ports)
	if out[0].Labels == nil {
		t.Fatal("expected labels map to be initialised")
	}
}
