package pipeline

import (
	"context"
	"sync"
	"time"
)

// ThrottleConfig controls the throttle behaviour.
type ThrottleConfig struct {
	// Rate is the maximum number of lines allowed per Window.
	Rate int
	// Window is the rolling time window for the rate limit.
	Window time.Duration
	// BurstSize allows short bursts above Rate before throttling kicks in.
	BurstSize int
}

// DefaultThrottleConfig returns a sensible default configuration.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		Rate:      1000,
		Window:    time.Second,
		BurstSize: 100,
	}
}

// Throttle enforces a per-window line rate, dropping lines that exceed it.
type Throttle struct {
	cfg     ThrottleConfig
	mu      sync.Mutex
	count   int
	burst   int
	windowStart time.Time
}

// NewThrottle creates a Throttle using the supplied config.
// If cfg.Rate <= 0 or cfg.Window <= 0 the defaults are applied.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	if cfg.Rate <= 0 || cfg.Window <= 0 {
		cfg = DefaultThrottleConfig()
	}
	if cfg.BurstSize < 0 {
		cfg.BurstSize = 0
	}
	return &Throttle{
		cfg:         cfg,
		windowStart: time.Now(),
	}
}

// Allow returns true if the line should be forwarded, false if it should be
// dropped. It is safe to call from multiple goroutines.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if now.Sub(t.windowStart) >= t.cfg.Window {
		t.count = 0
		t.burst = 0
		t.windowStart = now
	}

	limit := t.cfg.Rate
	if t.burst < t.cfg.BurstSize {
		limit = t.cfg.Rate + (t.cfg.BurstSize - t.burst)
	}

	if t.count >= limit {
		return false
	}

	t.count++
	if t.count > t.cfg.Rate {
		t.burst++
	}
	return true
}

// AllowCtx is like Allow but returns false immediately when ctx is cancelled.
func (t *Throttle) AllowCtx(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return t.Allow()
	}
}

// Reset clears the current window counters.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count = 0
	t.burst = 0
	t.windowStart = time.Now()
}
