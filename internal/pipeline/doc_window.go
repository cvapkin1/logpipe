// Package pipeline provides core log processing primitives for logpipe.
//
// # Window
//
// Window implements a sliding time-based aggregation buffer. Log lines are
// accumulated within a configurable time window (WindowConfig.Size) or until
// the buffer reaches its capacity (WindowConfig.MaxLines), whichever comes
// first. When either condition is met, the buffered lines are returned as a
// batch and the window resets.
//
// Typical usage:
//
//	w, err := pipeline.NewWindow(pipeline.DefaultWindowConfig())
//	if err != nil { ... }
//
//	for _, line := range incoming {
//	    if batch, ok := w.Add(line); ok {
//	        process(batch)
//	    }
//	}
//	// Flush remaining lines on shutdown.
//	if tail := w.Close(); len(tail) > 0 {
//	    process(tail)
//	}
//
// Window is safe for concurrent use.
package pipeline
