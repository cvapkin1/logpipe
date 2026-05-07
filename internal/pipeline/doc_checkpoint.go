// Package pipeline provides the core log-processing primitives for logpipe.
//
// # Checkpoint
//
// The Checkpoint type gives logpipe resume-after-restart semantics by
// persisting the last successfully processed byte offset (or any monotonic
// cursor) for each named source to a JSON file on disk.
//
// Usage:
//
//	cp, err := pipeline.NewCheckpoint("/var/lib/logpipe/checkpoint.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Restore the previous position before starting a source.
//	offset := cp.Get("nginx-access")
//
//	// After each successful write to a sink, advance the cursor.
//	_ = cp.Set("nginx-access", offset+int64(len(line)))
//
// The file is created automatically on the first Set call. Concurrent
// calls to Set, Get, and Delete are safe; each mutating call flushes
// the entire map to disk atomically via os.WriteFile.
package pipeline
