// Package pipeline provides core log routing and processing primitives
// for the logpipe aggregator.
//
// # Deduplicator
//
// The Deduplicator suppresses repeated log lines that appear within a
// configurable sliding time window. This is useful when a noisy service
// emits the same error message at high frequency and downstream sinks
// should only receive the first occurrence per window.
//
// Usage:
//
//	 cfg := pipeline.DefaultDeduplicatorConfig()
//	 cfg.WindowSize = 10 * time.Second
//	 dedup := pipeline.NewDeduplicator(cfg)
//
//	 for _, line := range incomingLines {
//	     if !dedup.IsDuplicate(line) {
//	         sink.Write(line)
//	     }
//	 }
//
// When MaxTracked is reached the entry with the oldest last-seen timestamp
// is evicted to make room for the new message, bounding memory usage.
package pipeline
