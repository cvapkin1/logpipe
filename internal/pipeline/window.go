package pipeline

import (
	"errors"
	"sync"
	"time"
)

// WindowConfig holds configuration for a sliding time window aggregator.
type WindowConfig struct {
	// Size is the duration of the window.
	Size time.Duration
	// MaxLines is the maximum number of lines to buffer within the window.
	MaxLines int
}

// DefaultWindowConfig returns a WindowConfig with sensible defaults.
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{
		Size:     10 * time.Second,
		MaxLines: 1000,
	}
}

// Window aggregates log lines within a sliding time window and flushes
// them when the window expires or the buffer is full.
type Window struct {
	cfg    WindowConfig
	mu     sync.Mutex
	buf    []string
	start  time.Time
	closed bool
}

// NewWindow creates a new Window with the given config.
// Returns an error if the config is invalid.
func NewWindow(cfg WindowConfig) (*Window, error) {
	if cfg.Size <= 0 {
		return nil, errors.New("window: Size must be positive")
	}
	if cfg.MaxLines <= 0 {
		return nil, errors.New("window: MaxLines must be positive")
	}
	return &Window{
		cfg:   cfg,
		buf:   make([]string, 0, cfg.MaxLines),
		start: time.Now(),
	}, nil
}

// Add appends a line to the window. Returns (lines, true) and resets the
// window when it is full or expired; otherwise returns (nil, false).
func (w *Window) Add(line string) ([]string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil, false
	}

	w.buf = append(w.buf, line)
	expired := time.Since(w.start) >= w.cfg.Size
	full := len(w.buf) >= w.cfg.MaxLines

	if expired || full {
		return w.flush(), true
	}
	return nil, false
}

// Flush forces a flush of the current window contents regardless of expiry.
func (w *Window) Flush() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flush()
}

// flush is the internal unlocked flush. Must be called with w.mu held.
func (w *Window) flush() []string {
	if len(w.buf) == 0 {
		w.start = time.Now()
		return nil
	}
	out := make([]string, len(w.buf))
	copy(out, w.buf)
	w.buf = w.buf[:0]
	w.start = time.Now()
	return out
}

// Close marks the window as closed and returns any buffered lines.
func (w *Window) Close() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
	return w.flush()
}
