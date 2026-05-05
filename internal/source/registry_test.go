package source

import (
	"context"
	"testing"
)

// stubSource is a minimal Source for testing.
type stubSource struct{ name string }

func (s *stubSource) Name() string                                      { return s.name }
func (s *stubSource) Start(_ context.Context, _ chan<- string) error    { return nil }
func (s *stubSource) Close() error                                      { return nil }

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	s := &stubSource{name: "app"}
	if err := reg.Register(s); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, err := reg.Get("app")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name() != "app" {
		t.Fatalf("expected 'app', got %q", got.Name())
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := NewRegistry()
	s := &stubSource{name: "dup"}
	if err := reg.Register(s); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	if err := reg.Register(s); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestRegistry_RegisterNil(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(nil); err == nil {
		t.Fatal("expected error when registering nil")
	}
}

func TestRegistry_GetNotFound(t *testing.T) {
	reg := NewRegistry()
	if _, err := reg.Get("missing"); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := NewRegistry()
	for _, n := range []string{"a", "b", "c"} {
		if err := reg.Register(&stubSource{name: n}); err != nil {
			t.Fatalf("Register %q: %v", n, err)
		}
	}
	names := reg.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
}

func TestRegistry_CloseAll(t *testing.T) {
	reg := NewRegistry()
	for _, n := range []string{"x", "y"} {
		_ = reg.Register(&stubSource{name: n})
	}
	errs := reg.CloseAll()
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}
