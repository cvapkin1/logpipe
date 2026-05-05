package pipeline_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/pipeline"
)

// mockDest is a simple in-memory Destination used in tests.
type mockDest struct {
	name    string
	entries []pipeline.LogEntry
	failOn  int // if > 0, return error after this many writes
}

func (m *mockDest) Write(e pipeline.LogEntry) error {
	if m.failOn > 0 && len(m.entries) >= m.failOn {
		return errors.New("write failed")
	}
	m.entries = append(m.entries, e)
	return nil
}
func (m *mockDest) Name() string { return m.name }

func makeRule(t *testing.T, minLevel filter.Level, contains string) *filter.Rule {
	t.Helper()
	rule, err := filter.NewRule(minLevel, contains, "")
	if err != nil {
		t.Fatalf("NewRule: %v", err)
	}
	return rule
}

func TestRouter_AddRoute_NilRoute(t *testing.T) {
	r := pipeline.NewRouter()
	if err := r.AddRoute(nil); err == nil {
		t.Fatal("expected error for nil route")
	}
}

func TestRouter_AddRoute_NilDestination(t *testing.T) {
	r := pipeline.NewRouter()
	err := r.AddRoute(&pipeline.Route{})
	if err == nil {
		t.Fatal("expected error for nil destination")
	}
}

func TestRouter_Dispatch_NoRoutes(t *testing.T) {
	r := pipeline.NewRouter()
	entry := pipeline.LogEntry{Source: "app", Level: filter.LevelInfo, Message: "hello"}
	n, err := r.Dispatch(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 forwards, got %d", n)
	}
}

func TestRouter_Dispatch_MatchingRoute(t *testing.T) {
	r := pipeline.NewRouter()
	dest := &mockDest{name: "out"}
	route := &pipeline.Route{
		Rules:       []*filter.Rule{makeRule(t, filter.LevelWarn, "")},
		Destination: dest,
	}
	_ = r.AddRoute(route)

	entry := pipeline.LogEntry{Source: "svc", Level: filter.LevelError, Message: "boom"}
	n, err := r.Dispatch(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 forward, got %d", n)
	}
	if len(dest.entries) != 1 || dest.entries[0].Message != "boom" {
		t.Fatal("entry not written to destination")
	}
}

func TestRouter_Dispatch_NonMatchingRoute(t *testing.T) {
	r := pipeline.NewRouter()
	dest := &mockDest{name: "out"}
	_ = r.AddRoute(&pipeline.Route{
		Rules:       []*filter.Rule{makeRule(t, filter.LevelError, "")},
		Destination: dest,
	})

	entry := pipeline.LogEntry{Source: "svc", Level: filter.LevelInfo, Message: "info msg"}
	n, _ := r.Dispatch(entry)
	if n != 0 {
		t.Fatalf("expected 0 forwards, got %d", n)
	}
	if len(dest.entries) != 0 {
		t.Fatal("entry should not have been written")
	}
}

func TestRouter_Dispatch_DestinationError(t *testing.T) {
	r := pipeline.NewRouter()
	dest := &mockDest{name: "failing", failOn: 0}
	_ = r.AddRoute(&pipeline.Route{
		Destination: dest,
	})

	entry := pipeline.LogEntry{Level: filter.LevelInfo, Message: fmt.Sprintf("msg")}
	_, err := r.Dispatch(entry)
	if err == nil {
		t.Fatal("expected error from failing destination")
	}
}
