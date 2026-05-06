// Package pipeline provides the core log routing and processing primitives
// for logpipe.
//
// # Buffer
//
// Buffer is a thread-safe, fixed-capacity ring buffer used to decouple log
// producers (sources) from consumers (sinks). When a source emits lines
// faster than a sink can drain them the buffer absorbs the burst up to its
// configured capacity.
//
// # Backpressure
//
// PushWithBackpressure wraps Buffer.Push with a configurable retry loop so
// that transient full-buffer conditions do not immediately cause line drops.
// Once all retries are exhausted, or the supplied context is cancelled, the
// line is dropped and a warning is logged.
//
// Typical usage:
//
//	buf := pipeline.NewBuffer(512)
//	cfg := pipeline.DefaultBackpressureConfig()
//	for line := range sourceLines {
//		pipeline.PushWithBackpressure(ctx, buf, line, cfg)
//	}
package pipeline
