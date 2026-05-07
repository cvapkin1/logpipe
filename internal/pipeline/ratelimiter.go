package pipeline

import (
	"context"
	"sync"
	"time"
)

// RateLimiter enforces a maximum number of log lines per second
// for a given pipeline route or source.
type RateLimiter struct {
	mu       sync.Mutex
	rate     int           // max events per interval
	interval time.Duration // measurement window
	count    int
	window   time.Time
}

// RateLimiterConfig holds configuration for a RateLimiter.
type RateLimiterConfig struct {
	// Rate is the maximum number of events allowed per Interval.
	Rate int
	// Interval is the duration of the sliding window (default: 1s).
	Interval time.Duration
}

// DefaultRateLimiterConfig returns a sensible default: 1000 lines/sec.
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Rate:     1000,
		Interval: time.Second,
	}
}

// NewRateLimiter creates a RateLimiter from the given config.
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.Rate <= 0 {
		cfg.Rate = 1000
	}
	return &RateLimiter{
		rate:     cfg.Rate,
		interval: cfg.Interval,
		window:   time.Now(),
	}
}

// Allow reports whether an event should be allowed through.
// It resets the counter when the current window expires.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if now.Sub(r.window) >= r.interval {
		r.count = 0
		r.window = now
	}
	if r.count >= r.rate {
		return false
	}
	r.count++
	return true
}

// Wait blocks until an event is allowed or the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		if r.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.interval / time.Duration(r.rate+1)):
		}
	}
}
