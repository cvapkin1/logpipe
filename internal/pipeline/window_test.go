package pipeline

import (
	"testing"
	"time"
)

func TestNewWindow_Defaults(t *testing.T) {
	cfg := DefaultWindowConfig()
	w, err := NewWindow(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestNewWindow_InvalidSize(t *testing.T) {
	_, err := NewWindow(WindowConfig{Size: 0, MaxLines: 10})
	if err == nil {
		t.Fatal("expected error for zero Size")
	}
}

func TestNewWindow_InvalidMaxLines(t *testing.T) {
	_, err := NewWindow(WindowConfig{Size: time.Second, MaxLines: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxLines")
	}
}

func TestWindow_Add_NoFlushBeforeExpiry(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 100})
	lines, flushed := w.Add("hello")
	if flushed {
		t.Fatalf("expected no flush, got %v", lines)
	}
	if lines != nil {
		t.Fatalf("expected nil lines, got %v", lines)
	}
}

func TestWindow_Add_FlushWhenFull(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 3})
	w.Add("a")
	w.Add("b")
	lines, flushed := w.Add("c")
	if !flushed {
		t.Fatal("expected flush when full")
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestWindow_Add_FlushWhenExpired(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 1 * time.Millisecond, MaxLines: 100})
	w.Add("first")
	time.Sleep(5 * time.Millisecond)
	lines, flushed := w.Add("second")
	if !flushed {
		t.Fatal("expected flush after expiry")
	}
	if len(lines) < 1 {
		t.Fatalf("expected at least 1 line, got %d", len(lines))
	}
}

func TestWindow_Flush_Empty(t *testing.T) {
	w, _ := NewWindow(DefaultWindowConfig())
	lines := w.Flush()
	if lines != nil {
		t.Fatalf("expected nil on empty flush, got %v", lines)
	}
}

func TestWindow_Flush_ReturnsBuffered(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 50})
	w.Add("x")
	w.Add("y")
	lines := w.Flush()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestWindow_Close_ReturnsPending(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 50})
	w.Add("line1")
	w.Add("line2")
	lines := w.Close()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines on close, got %d", len(lines))
	}
}

func TestWindow_Add_AfterClose_Ignored(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 50})
	w.Close()
	lines, flushed := w.Add("late")
	if flushed || lines != nil {
		t.Fatal("expected no flush after close")
	}
}

func TestWindow_ResetAfterFlush(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 10 * time.Second, MaxLines: 2})
	w.Add("a")
	w.Add("b") // triggers flush
	w.Add("c")
	lines := w.Flush()
	if len(lines) != 1 || lines[0] != "c" {
		t.Fatalf("expected [c] after reset, got %v", lines)
	}
}
