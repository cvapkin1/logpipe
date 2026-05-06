package pipeline

import (
	"sync/atomic"
	"time"
)

// Metrics tracks runtime statistics for the pipeline.
type Metrics struct {
	MessagesReceived  atomic.Int64
	MessagesForwarded atomic.Int64
	MessagesDropped   atomic.Int64
	BytesProcessed    atomic.Int64
	startedAt         time.Time
}

// NewMetrics creates a new Metrics instance with the start time set to now.
func NewMetrics() *Metrics {
	return &Metrics{
		startedAt: time.Now(),
	}
}

// RecordReceived increments the received message counter by n.
func (m *Metrics) RecordReceived(n int64) {
	m.MessagesReceived.Add(n)
}

// RecordForwarded increments the forwarded message counter by n.
func (m *Metrics) RecordForwarded(n int64) {
	m.MessagesForwarded.Add(n)
}

// RecordDropped increments the dropped message counter by n.
func (m *Metrics) RecordDropped(n int64) {
	m.MessagesDropped.Add(n)
}

// RecordBytes increments the total bytes processed counter by n.
func (m *Metrics) RecordBytes(n int64) {
	m.BytesProcessed.Add(n)
}

// Snapshot returns a point-in-time copy of the current metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		MessagesReceived:  m.MessagesReceived.Load(),
		MessagesForwarded: m.MessagesForwarded.Load(),
		MessagesDropped:   m.MessagesDropped.Load(),
		BytesProcessed:    m.BytesProcessed.Load(),
		Uptime:            time.Since(m.startedAt),
	}
}

// MetricsSnapshot is an immutable point-in-time view of pipeline metrics.
type MetricsSnapshot struct {
	MessagesReceived  int64
	MessagesForwarded int64
	MessagesDropped   int64
	BytesProcessed    int64
	Uptime            time.Duration
}
