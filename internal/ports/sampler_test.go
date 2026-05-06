package ports

import (
	"testing"
	"time"
)

func makeSamplerPort(port uint16) Port {
	return Port{Port: port, Protocol: "tcp"}
}

func TestSamplerDisabledPassesAll(t *testing.T) {
	cfg := DefaultSamplerConfig()
	cfg.Enabled = false
	s := NewSampler(cfg)

	input := []Port{makeSamplerPort(80), makeSamplerPort(443), makeSamplerPort(8080)}
	out := s.Sample(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d ports, got %d", len(input), len(out))
	}
}

func TestSamplerReservoirLimit(t *testing.T) {
	cfg := DefaultSamplerConfig()
	cfg.Enabled = true
	cfg.ReservoirSize = 3
	cfg.Window = time.Minute
	s := NewSampler(cfg)

	input := make([]Port, 10)
	for i := range input {
		input[i] = makeSamplerPort(uint16(1000 + i))
	}
	out := s.Sample(input)
	if len(out) > cfg.ReservoirSize {
		t.Fatalf("expected at most %d ports, got %d", cfg.ReservoirSize, len(out))
	}
}

func TestSamplerWindowReset(t *testing.T) {
	cfg := DefaultSamplerConfig()
	cfg.Enabled = true
	cfg.ReservoirSize = 2
	cfg.Window = 50 * time.Millisecond

	now := time.Now()
	s := NewSampler(cfg)
	s.now = func() time.Time { return now }

	s.Sample([]Port{makeSamplerPort(80), makeSamplerPort(443)})

	// Advance past the window.
	s.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	out := s.Sample([]Port{makeSamplerPort(8080)})
	if len(out) != 1 {
		t.Fatalf("expected 1 port after window reset, got %d", len(out))
	}
}

func TestSamplerReset(t *testing.T) {
	cfg := DefaultSamplerConfig()
	cfg.Enabled = true
	cfg.ReservoirSize = 5
	cfg.Window = time.Minute
	s := NewSampler(cfg)

	s.Sample([]Port{makeSamplerPort(80), makeSamplerPort(443)})
	s.Reset()

	out := s.Sample([]Port{makeSamplerPort(9090)})
	if len(out) != 1 {
		t.Fatalf("expected 1 port after reset, got %d", len(out))
	}
}

func TestSamplerAccumulatesAcrossCalls(t *testing.T) {
	cfg := DefaultSamplerConfig()
	cfg.Enabled = true
	cfg.ReservoirSize = 4
	cfg.Window = time.Minute
	s := NewSampler(cfg)

	s.Sample([]Port{makeSamplerPort(80), makeSamplerPort(443)})
	out := s.Sample([]Port{makeSamplerPort(8080), makeSamplerPort(8443)})
	if len(out) != 4 {
		t.Fatalf("expected 4 accumulated ports, got %d", len(out))
	}
}
