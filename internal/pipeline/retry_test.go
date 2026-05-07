package pipeline

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(context.Background(), DefaultRetryConfig(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	var calls int32
	err := Do(context.Background(), DefaultRetryConfig(), func(_ context.Context) error {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 2.0}
	calls := 0
	err := Do(context.Background(), cfg, func(_ context.Context) error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Do(ctx, DefaultRetryConfig(), func(_ context.Context) error {
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ContextCancelledDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cfg := RetryConfig{MaxAttempts: 5, InitialDelay: 500 * time.Millisecond, MaxDelay: 2 * time.Second, Multiplier: 2.0}
	start := time.Now()
	err := Do(ctx, cfg, func(_ context.Context) error { return errTemp })
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if elapsed > 300*time.Millisecond {
		t.Fatalf("context cancellation took too long: %v", elapsed)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
