package pipeline

import (
	"sync"
	"testing"
)

func TestNewBuffer_DefaultCapacity(t *testing.T) {
	b := NewBuffer(0)
	if b.Cap() != 256 {
		t.Fatalf("expected default cap 256, got %d", b.Cap())
	}
}

func TestBuffer_PushPop(t *testing.T) {
	b := NewBuffer(4)
	if err := b.Push("line1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Push("line2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 2 {
		t.Fatalf("expected len 2, got %d", b.Len())
	}
	line, ok := b.Pop()
	if !ok || line != "line1" {
		t.Fatalf("expected line1, got %q ok=%v", line, ok)
	}
	line, ok = b.Pop()
	if !ok || line != "line2" {
		t.Fatalf("expected line2, got %q ok=%v", line, ok)
	}
}

func TestBuffer_Full(t *testing.T) {
	b := NewBuffer(2)
	_ = b.Push("a")
	_ = b.Push("b")
	if err := b.Push("c"); err != ErrBufferFull {
		t.Fatalf("expected ErrBufferFull, got %v", err)
	}
}

func TestBuffer_PopEmpty(t *testing.T) {
	b := NewBuffer(4)
	_, ok := b.Pop()
	if ok {
		t.Fatal("expected ok=false on empty buffer")
	}
}

func TestBuffer_WrapAround(t *testing.T) {
	b := NewBuffer(3)
	_ = b.Push("x")
	_ = b.Push("y")
	b.Pop()
	_ = b.Push("z")
	_ = b.Push("w")
	if b.Len() != 3 {
		t.Fatalf("expected len 3, got %d", b.Len())
	}
}

func TestBuffer_ConcurrentAccess(t *testing.T) {
	b := NewBuffer(1024)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = b.Push("concurrent")
		}()
	}
	wg.Wait()
	if b.Len() > 1024 {
		t.Fatalf("buffer exceeded capacity: %d", b.Len())
	}
}
