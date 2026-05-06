package ports

import (
	"sync"
	"time"
)

// SamplerConfig controls reservoir sampling of port events.
type SamplerConfig struct {
	// Enabled toggles sampling; when false all ports pass through.
	Enabled bool
	// ReservoirSize is the maximum number of ports retained per window.
	ReservoirSize int
	// Window is the duration after which the reservoir resets.
	Window time.Duration
}

// DefaultSamplerConfig returns a sensible default.
func DefaultSamplerConfig() SamplerConfig {
	return SamplerConfig{
		Enabled:       false,
		ReservoirSize: 100,
		Window:        time.Minute,
	}
}

// Sampler performs reservoir sampling over a sliding time window.
type Sampler struct {
	cfg       SamplerConfig
	mu        sync.Mutex
	reservoir []Port
	windowEnd time.Time
	now       func() time.Time
}

// NewSampler constructs a Sampler from cfg.
func NewSampler(cfg SamplerConfig) *Sampler {
	return &Sampler{
		cfg: cfg,
		now: time.Now,
	}
}

// Sample accepts ports and returns a sampled subset.
// When sampling is disabled the input slice is returned unchanged.
func (s *Sampler) Sample(ports []Port) []Port {
	if !s.cfg.Enabled {
		return ports
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if now.After(s.windowEnd) {
		s.reservoir = s.reservoir[:0]
		s.windowEnd = now.Add(s.cfg.Window)
	}

	for _, p := range ports {
		if len(s.reservoir) < s.cfg.ReservoirSize {
			s.reservoir = append(s.reservoir, p)
		}
		// Once the reservoir is full additional ports are silently dropped
		// within the current window — callers receive only the retained set.
	}

	out := make([]Port, len(s.reservoir))
	copy(out, s.reservoir)
	return out
}

// Reset clears the reservoir and resets the window, useful for testing.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reservoir = s.reservoir[:0]
	s.windowEnd = time.Time{}
}
