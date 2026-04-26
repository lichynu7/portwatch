package config

import (
	"runtime"
	"testing"
)

func TestDefaultEnricherConfig(t *testing.T) {
	cfg := DefaultEnricherConfig()
	if !cfg.ResolveServices {
		t.Error("expected ResolveServices to be true by default")
	}
	if !cfg.LookupProcesses {
		t.Error("expected LookupProcesses to be true by default")
	}
}

func TestEnricherConfigValidateBothDisabled(t *testing.T) {
	cfg := EnricherConfig{ResolveServices: false, LookupProcesses: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error with both disabled: %v", err)
	}
}

func TestEnricherConfigValidateServiceOnly(t *testing.T) {
	cfg := EnricherConfig{ResolveServices: true, LookupProcesses: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEnricherConfigValidateProcessLookupNonLinux(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("skipping non-Linux test on Linux")
	}
	cfg := EnricherConfig{ResolveServices: false, LookupProcesses: true}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for process lookup on non-Linux, got nil")
	}
}

func TestEnricherConfigValidateProcessLookupLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping Linux-only test")
	}
	cfg := EnricherConfig{ResolveServices: true, LookupProcesses: true}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error on Linux: %v", err)
	}
}

func TestEnricherConfigIsLinuxOverride(t *testing.T) {
	orig := isLinux
	t.Cleanup(func() { isLinux = orig })

	isLinux = func() bool { return false }
	cfg := EnricherConfig{LookupProcesses: true}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when isLinux returns false")
	}

	isLinux = func() bool { return true }
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error when isLinux returns true: %v", err)
	}
}
