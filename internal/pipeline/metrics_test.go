package pipeline

import (
	"testing"
	"time"
)

func TestNewMetrics_InitialValues(t *testing.T) {
	m := NewMetrics()
	snap := m.Snapshot()

	if snap.MessagesReceived != 0 {
		t.Errorf("expected 0 received, got %d", snap.MessagesReceived)
	}
	if snap.MessagesForwarded != 0 {
		t.Errorf("expected 0 forwarded, got %d", snap.MessagesForwarded)
	}
	if snap.MessagesDropped != 0 {
		t.Errorf("expected 0 dropped, got %d", snap.MessagesDropped)
	}
	if snap.BytesProcessed != 0 {
		t.Errorf("expected 0 bytes, got %d", snap.BytesProcessed)
	}
}

func TestMetrics_RecordReceived(t *testing.T) {
	m := NewMetrics()
	m.RecordReceived(5)
	m.RecordReceived(3)

	if got := m.Snapshot().MessagesReceived; got != 8 {
		t.Errorf("expected 8, got %d", got)
	}
}

func TestMetrics_RecordForwarded(t *testing.T) {
	m := NewMetrics()
	m.RecordForwarded(10)

	if got := m.Snapshot().MessagesForwarded; got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

func TestMetrics_RecordDropped(t *testing.T) {
	m := NewMetrics()
	m.RecordDropped(2)

	if got := m.Snapshot().MessagesDropped; got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestMetrics_RecordBytes(t *testing.T) {
	m := NewMetrics()
	m.RecordBytes(1024)
	m.RecordBytes(512)

	if got := m.Snapshot().BytesProcessed; got != 1536 {
		t.Errorf("expected 1536, got %d", got)
	}
}

func TestMetrics_Snapshot_Uptime(t *testing.T) {
	m := NewMetrics()
	time.Sleep(10 * time.Millisecond)
	snap := m.Snapshot()

	if snap.Uptime < 10*time.Millisecond {
		t.Errorf("expected uptime >= 10ms, got %v", snap.Uptime)
	}
}

func TestMetrics_Snapshot_IsImmutable(t *testing.T) {
	m := NewMetrics()
	m.RecordReceived(1)
	snap := m.Snapshot()

	m.RecordReceived(99)

	if snap.MessagesReceived != 1 {
		t.Errorf("snapshot should not reflect later changes, got %d", snap.MessagesReceived)
	}
}
