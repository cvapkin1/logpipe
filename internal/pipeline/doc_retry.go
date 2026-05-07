// Package pipeline provides core log-processing primitives for logpipe.
//
// # Retry
//
// The retry module implements an exponential back-off retry policy for
// transient delivery failures. It is intended to wrap sink write operations
// so that temporary network or I/O errors do not immediately drop log lines.
//
// Usage:
//
//	cfg := pipeline.DefaultRetryConfig()
//	err := pipeline.Do(ctx, cfg, func(ctx context.Context) error {
//		return sink.Write(ctx, line)
//	})
//
// If all attempts are exhausted without success, Do returns
// ErrMaxAttemptsReached. If the context is cancelled during a back-off sleep
// or before the first attempt, the context error is returned immediately.
package pipeline
