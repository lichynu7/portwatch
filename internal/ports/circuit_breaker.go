package ports

import (
	"fmt"
	"sync"
	"time"
)

// CBState represents the state of a circuit breaker.
type CBState int

const (
	CBClosed CBState = iota // normal operation
	CBOpen                  // failing, reject calls
	CBHalfOpen             // testing recovery
)

func (s CBState) String() string {
	switch s {
	case CBClosed:
		return "closed"
	case CBOpen:
		return "open"
	case CBHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds tuning parameters for the breaker.
type CircuitBreakerConfig struct {
	FailureThreshold int           // consecutive failures before opening
	SuccessThreshold int           // consecutive successes in half-open before closing
	OpenTimeout      time.Duration // how long to stay open before probing
}

// DefaultCircuitBreakerConfig returns sensible defaults.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenTimeout:      30 * time.Second,
	}
}

// CircuitBreaker guards a downstream call (e.g. an alert handler) and
// stops hammering it when it is consistently failing.
type CircuitBreaker struct {
	cfg            CircuitBreakerConfig
	mu             sync.Mutex
	state          CBState
	failures       int
	successes      int
	lastTransition time.Time
}

// NewCircuitBreaker constructs a CircuitBreaker with the given config.
func NewCircuitBreaker(cfg CircuitBreakerConfig) (*CircuitBreaker, error) {
	if cfg.FailureThreshold <= 0 {
		return nil, fmt.Errorf("circuit breaker: FailureThreshold must be > 0")
	}
	if cfg.SuccessThreshold <= 0 {
		return nil, fmt.Errorf("circuit breaker: SuccessThreshold must be > 0")
	}
	if cfg.OpenTimeout <= 0 {
		return nil, fmt.Errorf("circuit breaker: OpenTimeout must be > 0")
	}
	return &CircuitBreaker{cfg: cfg, state: CBClosed}, nil
}

// Allow reports whether the caller is permitted to attempt the guarded call.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case CBClosed:
		return true
	case CBOpen:
		if time.Since(cb.lastTransition) >= cb.cfg.OpenTimeout {
			cb.state = CBHalfOpen
			cb.successes = 0
			return true
		}
		return false
	case CBHalfOpen:
		return true
	}
	return false
}

// RecordSuccess notifies the breaker of a successful call.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	if cb.state == CBHalfOpen {
		cb.successes++
		if cb.successes >= cb.cfg.SuccessThreshold {
			cb.state = CBClosed
			cb.lastTransition = time.Now()
		}
	}
}

// RecordFailure notifies the breaker of a failed call.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.successes = 0
	cb.failures++
	if cb.state == CBHalfOpen || cb.failures >= cb.cfg.FailureThreshold {
		cb.state = CBOpen
		cb.lastTransition = time.Now()
		cb.failures = 0
	}
}

// State returns the current breaker state (safe for concurrent use).
func (cb *CircuitBreaker) State() CBState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
