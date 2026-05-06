package ports

import (
	"testing"
	"time"
)

func defaultCBCfg() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      50 * time.Millisecond,
	}
}

func TestCircuitBreakerInitiallyAllows(t *testing.T) {
	cb, err := NewCircuitBreaker(defaultCBCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cb.Allow() {
		t.Fatal("expected Allow() == true in closed state")
	}
	if cb.State() != CBClosed {
		t.Fatalf("expected closed, got %s", cb.State())
	}
}

func TestCircuitBreakerOpensAfterThreshold(t *testing.T) {
	cb, _ := NewCircuitBreaker(defaultCBCfg())
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != CBOpen {
		t.Fatalf("expected open, got %s", cb.State())
	}
	if cb.Allow() {
		t.Fatal("expected Allow() == false when open")
	}
}

func TestCircuitBreakerTransitionsToHalfOpen(t *testing.T) {
	cb, _ := NewCircuitBreaker(defaultCBCfg())
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow() {
		t.Fatal("expected Allow() == true after open timeout")
	}
	if cb.State() != CBHalfOpen {
		t.Fatalf("expected half-open, got %s", cb.State())
	}
}

func TestCircuitBreakerClosesAfterSuccesses(t *testing.T) {
	cb, _ := NewCircuitBreaker(defaultCBCfg())
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	cb.Allow() // probe — transitions to half-open
	cb.RecordSuccess()
	cb.RecordSuccess()
	if cb.State() != CBClosed {
		t.Fatalf("expected closed after %d successes, got %s", 2, cb.State())
	}
}

func TestCircuitBreakerReopensFromHalfOpen(t *testing.T) {
	cb, _ := NewCircuitBreaker(defaultCBCfg())
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	cb.Allow()
	cb.RecordFailure() // fail during probe
	if cb.State() != CBOpen {
		t.Fatalf("expected open after failure in half-open, got %s", cb.State())
	}
}

func TestCircuitBreakerInvalidConfig(t *testing.T) {
	cases := []CircuitBreakerConfig{
		{FailureThreshold: 0, SuccessThreshold: 1, OpenTimeout: time.Second},
		{FailureThreshold: 1, SuccessThreshold: 0, OpenTimeout: time.Second},
		{FailureThreshold: 1, SuccessThreshold: 1, OpenTimeout: 0},
	}
	for _, cfg := range cases {
		_, err := NewCircuitBreaker(cfg)
		if err == nil {
			t.Errorf("expected error for config %+v", cfg)
		}
	}
}

func TestCBStateString(t *testing.T) {
	if CBClosed.String() != "closed" {
		t.Errorf("unexpected: %s", CBClosed.String())
	}
	if CBOpen.String() != "open" {
		t.Errorf("unexpected: %s", CBOpen.String())
	}
	if CBHalfOpen.String() != "half-open" {
		t.Errorf("unexpected: %s", CBHalfOpen.String())
	}
}
