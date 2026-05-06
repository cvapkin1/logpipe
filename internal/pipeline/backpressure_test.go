package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestPushWithBackpressure_Success(t *testing.T) {
	buf := NewBuffer(4)
	cfg := DefaultBackpressureConfig()
	ok := PushWithBackpressure(context.Background(), buf, "hello", cfg)
	if !ok {
		t.Fatal("expected push to succeed")
	}
	if buf.Len() != 1 {
		t.Fatalf("expected len 1, got %d", buf.Len())
	}
}

func TestPushWithBackpressure_DropsWhenFull(t *testing.T) {
	buf := NewBuffer(1)
	_ = buf.Push("blocker")
	cfg := BackpressureConfig{MaxRetries: 2, RetryInterval: time.Millisecond}
	ok := PushWithBackpressure(context.Background(), buf, "dropped", cfg)
	if ok {
		t.Fatal("expected push to fail (buffer full)")
	}
}

func TestPushWithBackpressure_SucceedsAfterRetry(t *testing.T) {
	buf := NewBuffer(1)
	_ = buf.Push("blocker")
	cfg := BackpressureConfig{MaxRetries: 5, RetryInterval: 5 * time.Millisecond}
	// Free up space after a short delay.
	go func() {
		time.Sleep(8 * time.Millisecond)
		buf.Pop()
	}()
	ok := PushWithBackpressure(context.Background(), buf, "retry-line", cfg)
	if !ok {
		t.Fatal("expected push to succeed after retry")
	}
}

func TestPushWithBackpressure_ContextCancelled(t *testing.T) {
	buf := NewBuffer(1)
	_ = buf.Push("blocker")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	cfg := BackpressureConfig{MaxRetries: 10, RetryInterval: 20 * time.Millisecond}
	ok := PushWithBackpressure(ctx, buf, "ctx-line", cfg)
	if ok {
		t.Fatal("expected push to fail due to cancelled context")
	}
}
