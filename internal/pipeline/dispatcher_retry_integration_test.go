package pipeline

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// flakyWriter simulates a sink that fails the first N writes then succeeds.
type flakyWriter struct {
	failFor int32
	calls   int32
	wrote   []string
}

func (f *flakyWriter) Write(line string) error {
	atomic.AddInt32(&f.calls, 1)
	if atomic.LoadInt32(&f.failFor) > 0 {
		atomic.AddInt32(&f.failFor, -1)
		return errors.New("sink unavailable")
	}
	f.wrote = append(f.wrote, line)
	return nil
}

func TestRetry_IntegrationWithDispatcher_EventualSuccess(t *testing.T) {
	sink := &flakyWriter{failFor: 2}
	cfg := RetryConfig{
		MaxAttempts:  5,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	line := "hello from retry integration"
	err := Do(context.Background(), cfg, func(_ context.Context) error {
		return sink.Write(line)
	})

	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if len(sink.wrote) != 1 || sink.wrote[0] != line {
		t.Fatalf("expected line to be written, got %v", sink.wrote)
	}
	if sink.calls != 3 {
		t.Fatalf("expected 3 total calls (2 failures + 1 success), got %d", sink.calls)
	}
}

func TestRetry_IntegrationWithDispatcher_ExhaustedDropsLine(t *testing.T) {
	sink := &flakyWriter{failFor: 10}
	cfg := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}

	err := Do(context.Background(), cfg, func(_ context.Context) error {
		return sink.Write("dropped line")
	})

	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if len(sink.wrote) != 0 {
		t.Fatalf("expected no lines written, got %v", sink.wrote)
	}
}
