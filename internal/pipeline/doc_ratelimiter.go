// Package pipeline provides core log routing, buffering, and flow-control
// primitives for the logpipe aggregator.
//
// # Rate Limiter
//
// RateLimiter enforces a maximum throughput (events per interval) on any
// pipeline stage. It uses a fixed-window counter that resets at the start
// of each interval.
//
// Basic usage:
//
//	rl := pipeline.NewRateLimiter(pipeline.RateLimiterConfig{
//		Rate:     500,
//		Interval: time.Second,
//	})
//
//	// Non-blocking check:
//	if rl.Allow() {
//		// forward the log line
//	}
//
//	// Blocking check (respects context cancellation):
//	if err := rl.Wait(ctx); err != nil {
//		// context cancelled or deadline exceeded
//	}
//
// Use DefaultRateLimiterConfig() for a 1000 lines/sec default.
package pipeline
