package pipeline

import (
	"sync"
	"time"
)

// DeduplicatorConfig holds configuration for the log deduplicator.
type DeduplicatorConfig struct {
	// WindowSize is the duration within which duplicate messages are suppressed.
	WindowSize time.Duration
	// MaxTracked is the maximum number of unique messages to track.
	MaxTracked int
}

// DefaultDeduplicatorConfig returns a DeduplicatorConfig with sensible defaults.
func DefaultDeduplicatorConfig() DeduplicatorConfig {
	return DeduplicatorConfig{
		WindowSize: 5 * time.Second,
		MaxTracked: 1024,
	}
}

// entry holds metadata about a seen message.
type entry struct {
	count    int
	firstAt  time.Time
	lastAt   time.Time
}

// Deduplicator suppresses repeated log lines within a sliding time window.
type Deduplicator struct {
	mu     sync.Mutex
	cfg    DeduplicatorConfig
	seen   map[string]*entry
}

// NewDeduplicator creates a Deduplicator with the given config.
func NewDeduplicator(cfg DeduplicatorConfig) *Deduplicator {
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = DefaultDeduplicatorConfig().WindowSize
	}
	if cfg.MaxTracked <= 0 {
		cfg.MaxTracked = DefaultDeduplicatorConfig().MaxTracked
	}
	return &Deduplicator{
		cfg:  cfg,
		seen: make(map[string]*entry, cfg.MaxTracked),
	}
}

// IsDuplicate returns true if the message was seen within the current window.
// It records the message as seen if it is new or the window has expired.
func (d *Deduplicator) IsDuplicate(msg string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	if e, ok := d.seen[msg]; ok {
		if now.Sub(e.firstAt) < d.cfg.WindowSize {
			e.count++
			e.lastAt = now
			return true
		}
		// Window expired — reset entry.
		e.count = 1
		e.firstAt = now
		e.lastAt = now
		return false
	}

	// Evict oldest entry if at capacity.
	if len(d.seen) >= d.cfg.MaxTracked {
		d.evictOldest(now)
	}

	d.seen[msg] = &entry{count: 1, firstAt: now, lastAt: now}
	return false
}

// evictOldest removes the entry with the oldest lastAt timestamp.
func (d *Deduplicator) evictOldest(now time.Time) {
	var oldest string
	var oldestTime time.Time
	for k, e := range d.seen {
		if oldest == "" || e.lastAt.Before(oldestTime) {
			oldest = k
			oldestTime = e.lastAt
		}
	}
	delete(d.seen, oldest)
}

// Reset clears all tracked messages.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]*entry, d.cfg.MaxTracked)
}
