package pipeline

import (
	"testing"
	"time"
)

func TestCircuitBreaker_InitiallyClosed(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	if cb.State() != "closed" {
		t.Fatalf("expected closed, got %s", cb.State())
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 3
	cb := NewCircuitBreaker(cfg)

	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != "open" {
		t.Fatalf("expected open, got %s", cb.State())
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1
	cfg.OpenDuration = 10 * time.Millisecond
	cb := NewCircuitBreaker(cfg)

	cb.RecordFailure()
	if cb.State() != "open" {
		t.Fatalf("expected open")
	}

	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil after open duration, got %v", err)
	}
	if cb.State() != "half-open" {
		t.Fatalf("expected half-open, got %s", cb.State())
	}
}

func TestCircuitBreaker_ClosesAfterSuccessThreshold(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1
	cfg.SuccessThreshold = 2
	cfg.OpenDuration = 10 * time.Millisecond
	cb := NewCircuitBreaker(cfg)

	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow() // transitions to half-open

	cb.RecordSuccess()
	if cb.State() != "half-open" {
		t.Fatalf("expected still half-open after 1 success")
	}
	cb.RecordSuccess()
	if cb.State() != "closed" {
		t.Fatalf("expected closed after success threshold, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1
	cfg.OpenDuration = 10 * time.Millisecond
	cb := NewCircuitBreaker(cfg)

	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow()

	cb.RecordFailure()
	if cb.State() != "open" {
		t.Fatalf("expected open after failure in half-open, got %s", cb.State())
	}
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 3
	cb := NewCircuitBreaker(cfg)

	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	cb.RecordFailure()
	cb.RecordFailure()
	// Only 2 failures since last success — should still be closed
	if cb.State() != "closed" {
		t.Fatalf("expected closed, got %s", cb.State())
	}
}

func TestNewCircuitBreaker_DefaultsApplied(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{})
	def := DefaultCircuitBreakerConfig()
	if cb.cfg.FailureThreshold != def.FailureThreshold {
		t.Errorf("expected default FailureThreshold %d, got %d", def.FailureThreshold, cb.cfg.FailureThreshold)
	}
	if cb.cfg.OpenDuration != def.OpenDuration {
		t.Errorf("expected default OpenDuration %v, got %v", def.OpenDuration, cb.cfg.OpenDuration)
	}
}
