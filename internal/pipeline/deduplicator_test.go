package pipeline

import (
	"testing"
	"time"
)

func TestDeduplicator_NewMessage_NotDuplicate(t *testing.T) {
	d := NewDeduplicator(DefaultDeduplicatorConfig())
	if d.IsDuplicate("hello world") {
		t.Error("expected first occurrence to not be a duplicate")
	}
}

func TestDeduplicator_RepeatedMessage_IsDuplicate(t *testing.T) {
	d := NewDeduplicator(DefaultDeduplicatorConfig())
	d.IsDuplicate("repeated line")
	if !d.IsDuplicate("repeated line") {
		t.Error("expected second occurrence to be a duplicate")
	}
}

func TestDeduplicator_WindowExpiry_NotDuplicate(t *testing.T) {
	cfg := DeduplicatorConfig{WindowSize: 50 * time.Millisecond, MaxTracked: 64}
	d := NewDeduplicator(cfg)
	d.IsDuplicate("expiring line")
	time.Sleep(80 * time.Millisecond)
	if d.IsDuplicate("expiring line") {
		t.Error("expected message to not be duplicate after window expiry")
	}
}

func TestDeduplicator_DifferentMessages_NotDuplicate(t *testing.T) {
	d := NewDeduplicator(DefaultDeduplicatorConfig())
	d.IsDuplicate("line one")
	if d.IsDuplicate("line two") {
		t.Error("different messages should not be considered duplicates")
	}
}

func TestDeduplicator_Reset_ClearsState(t *testing.T) {
	d := NewDeduplicator(DefaultDeduplicatorConfig())
	d.IsDuplicate("some message")
	d.Reset()
	if d.IsDuplicate("some message") {
		t.Error("expected message to not be duplicate after Reset")
	}
}

func TestDeduplicator_MaxTracked_Eviction(t *testing.T) {
	cfg := DeduplicatorConfig{WindowSize: 5 * time.Second, MaxTracked: 3}
	d := NewDeduplicator(cfg)

	d.IsDuplicate("msg-a")
	d.IsDuplicate("msg-b")
	d.IsDuplicate("msg-c")
	// Adding a 4th should evict the oldest.
	d.IsDuplicate("msg-d")

	if len(d.seen) > 3 {
		t.Errorf("expected at most 3 tracked entries, got %d", len(d.seen))
	}
}

func TestDeduplicator_DefaultConfig_Fallbacks(t *testing.T) {
	d := NewDeduplicator(DeduplicatorConfig{})
	if d.cfg.WindowSize != DefaultDeduplicatorConfig().WindowSize {
		t.Errorf("expected default window size, got %v", d.cfg.WindowSize)
	}
	if d.cfg.MaxTracked != DefaultDeduplicatorConfig().MaxTracked {
		t.Errorf("expected default max tracked, got %d", d.cfg.MaxTracked)
	}
}
