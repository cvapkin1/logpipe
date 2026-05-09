package pipeline

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is in the open state.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// circuitState represents the state of the circuit breaker.
type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

// CircuitBreakerConfig holds configuration for the CircuitBreaker.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures before opening.
	FailureThreshold int
	// SuccessThreshold is the number of consecutive successes to close from half-open.
	SuccessThreshold int
	// OpenDuration is how long the circuit stays open before moving to half-open.
	OpenDuration time.Duration
}

// DefaultCircuitBreakerConfig returns a CircuitBreakerConfig with sensible defaults.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenDuration:     10 * time.Second,
	}
}

// CircuitBreaker guards a downstream sink or operation against cascading failures.
type CircuitBreaker struct {
	cfg            CircuitBreakerConfig
	mu             sync.Mutex
	state          circuitState
	failureCount   int
	successCount   int
	openedAt       time.Time
}

// NewCircuitBreaker creates a CircuitBreaker with the given config.
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = DefaultCircuitBreakerConfig().FailureThreshold
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = DefaultCircuitBreakerConfig().SuccessThreshold
	}
	if cfg.OpenDuration <= 0 {
		cfg.OpenDuration = DefaultCircuitBreakerConfig().OpenDuration
	}
	return &CircuitBreaker{cfg: cfg, state: stateClosed}
}

// Allow reports whether the operation should proceed.
// Returns ErrCircuitOpen if the circuit is open.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case stateOpen:
		if time.Since(cb.openedAt) >= cb.cfg.OpenDuration {
			cb.state = stateHalfOpen
			cb.successCount = 0
			return nil
		}
		return ErrCircuitOpen
	default:
		return nil
	}
}

// RecordSuccess records a successful operation.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount = 0
	if cb.state == stateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.cfg.SuccessThreshold {
			cb.state = stateClosed
			cb.successCount = 0
		}
	}
}

// RecordFailure records a failed operation and may open the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.successCount = 0
	cb.failureCount++
	if cb.state == stateHalfOpen || cb.failureCount >= cb.cfg.FailureThreshold {
		cb.state = stateOpen
		cb.openedAt = time.Now()
		cb.failureCount = 0
	}
}

// State returns the current circuit state as a string.
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case stateOpen:
		return "open"
	case stateHalfOpen:
		return "half-open"
	default:
		return "closed"
	}
}
