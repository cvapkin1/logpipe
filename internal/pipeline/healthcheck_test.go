package pipeline

import (
	"testing"
	"time"
)

func TestHealthChecker_Healthy(t *testing.T) {
	m := NewMetrics()
	m.RecordReceived()
	m.RecordReceived()
	m.RecordForwarded()
	m.RecordForwarded()

	hc := NewHealthChecker(m, 0.1, 30*time.Second)
	report := hc.Check()

	if report.Status != StatusHealthy {
		t.Errorf("expected Healthy, got %d: %s", report.Status, report.Message)
	}
	if report.Message != "ok" {
		t.Errorf("expected message 'ok', got %q", report.Message)
	}
}

func TestHealthChecker_Degraded_HighDropRate(t *testing.T) {
	m := NewMetrics()
	for i := 0; i < 5; i++ {
		m.RecordReceived()
	}
	for i := 0; i < 4; i++ {
		m.RecordDropped()
	}

	hc := NewHealthChecker(m, 0.1, 30*time.Second)
	report := hc.Check()

	if report.Status != StatusDegraded {
		t.Errorf("expected Degraded, got %d", report.Status)
	}
	if report.DropRate < 0.1 {
		t.Errorf("expected drop rate >= 0.1, got %f", report.DropRate)
	}
}

func TestHealthChecker_Unhealthy_Stale(t *testing.T) {
	m := NewMetrics()
	m.RecordReceived() // triggers LastReceivedAt to be set in the past

	// Use a very short stale window so the single past message is already stale.
	hc := NewHealthChecker(m, 0.1, 1*time.Nanosecond)
	time.Sleep(2 * time.Millisecond)
	report := hc.Check()

	if report.Status != StatusUnhealthy {
		t.Errorf("expected Unhealthy, got %d: %s", report.Status, report.Message)
	}
}

func TestHealthChecker_DefaultThresholds(t *testing.T) {
	m := NewMetrics()
	hc := NewHealthChecker(m, 0, 0) // zero values → defaults applied

	if hc.degradedDropRate != 0.1 {
		t.Errorf("expected default degradedDropRate 0.1, got %f", hc.degradedDropRate)
	}
	if hc.staleDuration != 30*time.Second {
		t.Errorf("expected default staleDuration 30s, got %v", hc.staleDuration)
	}
}

func TestHealthChecker_CheckedAt_Recent(t *testing.T) {
	m := NewMetrics()
	hc := NewHealthChecker(m, 0.1, 30*time.Second)
	before := time.Now()
	report := hc.Check()
	after := time.Now()

	if report.CheckedAt.Before(before) || report.CheckedAt.After(after) {
		t.Errorf("CheckedAt %v not between %v and %v", report.CheckedAt, before, after)
	}
}
