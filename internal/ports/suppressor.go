package ports

import (
	"sync"
	"time"
)

// SuppressorConfig holds configuration for the alert suppressor.
type SuppressorConfig struct {
	// QuietHoursStart is the hour (0-23) when suppression begins.
	QuietHoursStart int
	// QuietHoursEnd is the hour (0-23) when suppression ends.
	QuietHoursEnd int
	// MinSeverity is the minimum severity level that bypasses quiet hours.
	MinSeverity Severity
}

// DefaultSuppressorConfig returns a SuppressorConfig with quiet hours disabled.
func DefaultSuppressorConfig() SuppressorConfig {
	return SuppressorConfig{
		QuietHoursStart: -1,
		QuietHoursEnd:   -1,
		MinSeverity:     SeverityCritical,
	}
}

// Suppressor decides whether an alert should be suppressed based on
// time-of-day quiet hours and severity overrides.
type Suppressor struct {
	cfg  SuppressorConfig
	now  func() time.Time
	mu   sync.Mutex
}

// NewSuppressor creates a Suppressor with the given config.
func NewSuppressor(cfg SuppressorConfig) *Suppressor {
	return &Suppressor{
		cfg: cfg,
		now: time.Now,
	}
}

// Suppress returns true if the alert for the given port should be suppressed.
func (s *Suppressor) Suppress(p Port) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cfg.QuietHoursStart < 0 || s.cfg.QuietHoursEnd < 0 {
		return false
	}

	// Critical (or above min) alerts always pass through.
	if p.Severity >= s.cfg.MinSeverity {
		return false
	}

	hour := s.now().Hour()
	return s.inQuietWindow(hour)
}

// inQuietWindow checks whether hour falls within the configured quiet range.
// Handles overnight windows (e.g. 22 -> 06).
func (s *Suppressor) inQuietWindow(hour int) bool {
	start := s.cfg.QuietHoursStart
	end := s.cfg.QuietHoursEnd
	if start <= end {
		return hour >= start && hour < end
	}
	// Overnight window.
	return hour >= start || hour < end
}

// ApplySuppressor returns a filter function that drops ports suppressed by s.
func ApplySuppressor(s *Suppressor) func([]Port) []Port {
	return func(ports []Port) []Port {
		if s == nil {
			return ports
		}
		out := ports[:0]
		for _, p := range ports {
			if !s.Suppress(p) {
				out = append(out, p)
			}
		}
		return out
	}
}
