package pipeline

import (
	"context"
	"errors"
	"time"
)

// RetryConfig holds configuration for the retry policy.
type RetryConfig struct {
	// MaxAttempts is the maximum number of delivery attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the first retry.
	InitialDelay time.Duration
	// MaxDelay caps the exponential back-off delay.
	MaxDelay time.Duration
	// Multiplier is the factor applied to the delay after each attempt.
	Multiplier float64
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryFunc is the signature of the operation that will be retried.
type RetryFunc func(ctx context.Context) error

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Do executes fn according to the retry policy defined in cfg.
// It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, cfg RetryConfig, fn RetryFunc) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 1.0 {
		cfg.Multiplier = 2.0
	}

	delay := cfg.InitialDelay
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}

		if attempt < cfg.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}

	return ErrMaxAttemptsReached
}
