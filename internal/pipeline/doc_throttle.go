// Package pipeline provides the core log-processing primitives for logpipe.
//
// # Throttle
//
// Throttle enforces a maximum line throughput over a rolling time window,
// optionally allowing short-lived bursts above the base rate.
//
// Basic usage:
//
//	th := pipeline.NewThrottle(pipeline.ThrottleConfig{
//		Rate:      500,           // max 500 lines per window
//		Window:    time.Second,
//		BurstSize: 50,            // allow up to 50 extra lines in a burst
//	})
//
//	for _, line := range lines {
//		if th.AllowCtx(ctx) {
//			forward(line)
//		} else {
//			metrics.RecordDropped(1)
//		}
//	}
//
// Call Reset to clear counters between test cases or pipeline restarts.
// All methods are safe for concurrent use.
package pipeline
