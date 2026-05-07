package pipeline

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// mockDedupSink records every line written to it.
type mockDedupSink struct {
	mu    sync.Mutex
	lines []string
}

func (s *mockDedupSink) Write(line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lines = append(s.lines, line)
	return nil
}

func (s *mockDedupSink) Lines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.lines))
	copy(out, s.lines)
	return out
}

// dedupFilter wraps a Deduplicator to drop duplicate lines before forwarding.
func dedupFilter(d *Deduplicator, sink *mockDedupSink, line string) {
	if !d.IsDuplicate(line) {
		_ = sink.Write(line)
	}
}

func TestDeduplicator_IntegrationWithSink_SuppressDuplicates(t *testing.T) {
	cfg := DeduplicatorConfig{WindowSize: 2 * time.Second, MaxTracked: 128}
	d := NewDeduplicator(cfg)
	sink := &mockDedupSink{}

	lines := []string{
		"error: connection refused",
		"error: connection refused",
		"info: server started",
		"error: connection refused",
		"info: server started",
	}

	for _, l := range lines {
		dedupFilter(d, sink, l)
	}

	got := sink.Lines()
	if len(got) != 2 {
		t.Fatalf("expected 2 unique lines forwarded, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0], "connection refused") {
		t.Errorf("unexpected first line: %s", got[0])
	}
	if !strings.Contains(got[1], "server started") {
		t.Errorf("unexpected second line: %s", got[1])
	}
}

func TestDeduplicator_IntegrationWithSink_AllowsAfterWindowExpiry(t *testing.T) {
	cfg := DeduplicatorConfig{WindowSize: 50 * time.Millisecond, MaxTracked: 128}
	d := NewDeduplicator(cfg)
	sink := &mockDedupSink{}

	dedupFilter(d, sink, "warn: disk full")
	dedupFilter(d, sink, "warn: disk full") // duplicate — suppressed
	time.Sleep(80 * time.Millisecond)
	dedupFilter(d, sink, "warn: disk full") // window expired — forwarded

	got := sink.Lines()
	if len(got) != 2 {
		t.Fatalf("expected 2 forwarded lines after window expiry, got %d", len(got))
	}
}
