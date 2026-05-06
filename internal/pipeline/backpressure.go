package pipeline

import (
	"context"
	"log"
	"time"
)

// BackpressureConfig controls retry behaviour when a buffer is full.
type BackpressureConfig struct {
	// MaxRetries is the number of push attempts before dropping the line.
	MaxRetries int
	// RetryInterval is the wait between retries.
	RetryInterval time.Duration
}

// DefaultBackpressureConfig returns sensible defaults.
func DefaultBackpressureConfig() BackpressureConfig {
	return BackpressureConfig{
		MaxRetries:    5,
		RetryInterval: 10 * time.Millisecond,
	}
}

// PushWithBackpressure attempts to push a line into buf, retrying up to
// cfg.MaxRetries times. Returns true if the line was accepted, false if
// dropped. The context is checked between retries.
func PushWithBackpressure(ctx context.Context, buf *Buffer, line string, cfg BackpressureConfig) bool {
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if err := buf.Push(line); err == nil {
			return true
		}
		if attempt == cfg.MaxRetries {
			break
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(cfg.RetryInterval):
		}
	}
	log.Printf("pipeline: dropped line after %d retries: %.80q", cfg.MaxRetries, line)
	return false
}
