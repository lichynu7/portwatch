package ports

import (
	"sync"
	"time"
)

// ThrottleConfig controls per-port alert throttling behaviour.
type ThrottleConfig struct {
	// Window is the minimum duration between repeated alerts for the same port.
	Window time.Duration
	// MaxBurst is the number of alerts allowed before throttling kicks in.
	MaxBurst int
}

// DefaultThrottleConfig returns a sensible default throttle configuration.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		Window:   5 * time.Minute,
		MaxBurst: 1,
	}
}

// portThrottleState tracks alert history for a single port key.
type portThrottleState struct {
	count     int
	windowEnd time.Time
}

// Throttle suppresses repeated alerts for the same port within a configurable
// time window, allowing a small burst before engaging throttling.
type Throttle struct {
	mu     sync.Mutex
	cfg    ThrottleConfig
	states map[string]*portThrottleState
	now    func() time.Time
}

// NewThrottle creates a Throttle with the given configuration.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	return &Throttle{
		cfg:    cfg,
		states: make(map[string]*portThrottleState),
		now:    time.Now,
	}
}

// Allow returns true if an alert for key should be delivered, and false if it
// should be suppressed. It advances the internal counter for key.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	s, ok := t.states[key]
	if !ok || now.After(s.windowEnd) {
		t.states[key] = &portThrottleState{
			count:     1,
			windowEnd: now.Add(t.cfg.Window),
		}
		return true
	}

	s.count++
	return s.count <= t.cfg.MaxBurst
}

// Reset clears throttle state for key, allowing the next alert through
// unconditionally.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.states, key)
}
