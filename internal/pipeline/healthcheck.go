package pipeline

import (
	"sync"
	"time"
)

// HealthStatus represents the health state of a pipeline component.
type HealthStatus int

const (
	StatusHealthy   HealthStatus = iota
	StatusDegraded               // high drop rate or backpressure
	StatusUnhealthy              // no messages processed recently
)

// HealthReport is a point-in-time snapshot of pipeline health.
type HealthReport struct {
	Status      HealthStatus
	CheckedAt   time.Time
	DropRate    float64 // drops / (received + 1)
	LastMessage time.Time
	Message     string
}

// HealthChecker evaluates pipeline Metrics and returns a HealthReport.
type HealthChecker struct {
	mu              sync.Mutex
	metrics         *Metrics
	degradedDropRate float64
	staleDuration   time.Duration
}

// NewHealthChecker creates a HealthChecker backed by the given Metrics.
// degradedDropRate is the fraction of dropped messages (0–1) above which
// the pipeline is considered degraded. staleDuration is the maximum time
// without a received message before the pipeline is considered unhealthy.
func NewHealthChecker(m *Metrics, degradedDropRate float64, staleDuration time.Duration) *HealthChecker {
	if degradedDropRate <= 0 {
		degradedDropRate = 0.1
	}
	if staleDuration <= 0 {
		staleDuration = 30 * time.Second
	}
	return &HealthChecker{
		metrics:         m,
		degradedDropRate: degradedDropRate,
		staleDuration:   staleDuration,
	}
}

// Check evaluates the current metrics and returns a HealthReport.
func (h *HealthChecker) Check() HealthReport {
	h.mu.Lock()
	defer h.mu.Unlock()

	snap := h.metrics.Snapshot()
	now := time.Now()

	report := HealthReport{
		Status:      StatusHealthy,
		CheckedAt:   now,
		LastMessage: snap.LastReceivedAt,
	}

	received := snap.Received
	dropRate := float64(snap.Dropped) / float64(received+1)
	report.DropRate = dropRate

	if !snap.LastReceivedAt.IsZero() && now.Sub(snap.LastReceivedAt) > h.staleDuration {
		report.Status = StatusUnhealthy
		report.Message = "no messages received within stale window"
		return report
	}

	if dropRate >= h.degradedDropRate {
		report.Status = StatusDegraded
		report.Message = "drop rate exceeds threshold"
		return report
	}

	report.Message = "ok"
	return report
}
