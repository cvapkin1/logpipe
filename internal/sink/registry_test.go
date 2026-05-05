package sink

import (
	"sort"
	"testing"
)

type mockSink struct {
	name   string
	closed bool
}

func (m *mockSink) Write(line []byte) error { return nil }
func (m *mockSink) Close() error           { m.closed = true; return nil }
func (m *mockSink) Name() string           { return m.name }

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	s := &mockSink{name: "s1"}

	if err := r.Register(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("s1")
	if !ok {
		t.Fatal("expected sink to be found")
	}
	if got.Name() != "s1" {
		t.Errorf("expected name 's1', got %q", got.Name())
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockSink{name: "dup"})
	err := r.Register(&mockSink{name: "dup"})
	if err == nil {
		t.Error("expected error for duplicate sink name")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Get("missing")
	if ok {
		t.Error("expected false for missing sink")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockSink{name: "b"})
	r.Register(&mockSink{name: "a"})
	r.Register(&mockSink{name: "c"})

	names := r.Names()
	sort.Strings(names)
	if len(names) != 3 || names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestRegistry_CloseAll(t *testing.T) {
	r := NewRegistry()
	s1 := &mockSink{name: "x"}
	s2 := &mockSink{name: "y"}
	r.Register(s1)
	r.Register(s2)

	if err := r.CloseAll(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s1.closed || !s2.closed {
		t.Error("expected both sinks to be closed")
	}
}
