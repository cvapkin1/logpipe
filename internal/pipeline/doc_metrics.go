// Package pipeline provides the core log routing and processing primitives
// used by logpipe.
//
// # Metrics
//
// The Metrics type offers lock-free, atomic counters for tracking pipeline
// throughput at runtime. Counters cover four dimensions:
//
//   - MessagesReceived  – total log lines ingested from all sources.
//   - MessagesForwarded – lines successfully written to at least one sink.
//   - MessagesDropped   – lines discarded due to back-pressure or filter rules.
//   - BytesProcessed    – cumulative byte volume of all ingested messages.
//
// Use Snapshot to obtain a consistent, point-in-time read of all counters
// without holding any lock:
//
//	snap := metrics.Snapshot()
//	fmt.Printf("received=%d forwarded=%d dropped=%d bytes=%d uptime=%s\n",
//		snap.MessagesReceived, snap.MessagesForwarded,
//		snap.MessagesDropped, snap.BytesProcessed, snap.Uptime)
//
// Metrics is safe for concurrent use by multiple goroutines.
package pipeline
