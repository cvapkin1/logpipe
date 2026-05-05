package pipeline

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockSource implements source.Source for testing.
type mockSource struct {
	name  string
	lines []string
}

func (m *mockSource) Name() string { return m.name }
func (m *mockSource) Start(_ context.Context) (<-chan string, error) {
	ch := make(chan string, len(m.lines))
	for _, l := range m.lines {
		ch <- l
	}
	close(ch)
	return ch, nil
}
func (m *mockSource) Close() error { return nil }

// collectSink implements sink.Sink for testing.
type collectSink struct {
	mu   sync.Mutex
	got  []string
	name string
}

func (c *collectSink) Name() string { return c.name }
func (c *collectSink) Write(line string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.got = append(c.got, line)
	return nil
}
func (c *collectSink) Close() error { return nil }

func TestDispatcher_RoutesLinesToSink(t *testing.T) {
	router := NewRouter()
	cs := &collectSink{name: "test"}
	// nil rule means accept-all in matchesAll
	if err := router.AddRoute(nil, cs); err != nil {
		t.Fatalf("AddRoute: %v", err)
	}

	src := &mockSource{name: "mock", lines: []string{"hello", "world"}}
	d := NewDispatcher(router)
	d.AddSource(src)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	d.Run(ctx)

	cs.mu.Lock()
	defer cs.mu.Unlock()
	if len(cs.got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(cs.got))
	}
	if !strings.Contains(cs.got[0], "hello") || !strings.Contains(cs.got[1], "world") {
		t.Errorf("unexpected lines: %v", cs.got)
	}
}

func TestDispatcher_NoSources(t *testing.T) {
	router := NewRouter()
	d := NewDispatcher(router)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Should return immediately with no sources.
	done := make(chan struct{})
	go func() { d.Run(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return promptly with no sources")
	}
}
