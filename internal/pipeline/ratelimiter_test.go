package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter_Defaults(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{})
	if rl.rate != 1000 {
		t.Errorf("expected default rate 1000, got %d", rl.rate)
	}
	if rl.interval != time.Second {
		t.Errorf("expected default interval 1s, got %v", rl.interval)
	}
}

func TestRateLimiter_AllowsUpToRate(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 5, Interval: time.Second})
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if rl.Allow() {
		t.Error("expected Allow()=false after rate exceeded")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 2, Interval: 50 * time.Millisecond})
	rl.Allow()
	rl.Allow()
	if rl.Allow() {
		t.Error("expected false before window resets")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow() {
		t.Error("expected true after window reset")
	}
}

func TestRateLimiter_DefaultConfig(t *testing.T) {
	cfg := DefaultRateLimiterConfig()
	if cfg.Rate != 1000 {
		t.Errorf("expected rate 1000, got %d", cfg.Rate)
	}
	if cfg.Interval != time.Second {
		t.Errorf("expected interval 1s, got %v", cfg.Interval)
	}
}

func TestRateLimiter_Wait_Success(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 10, Interval: time.Second})
	ctx := context.Background()
	if err := rl.Wait(ctx); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRateLimiter_Wait_ContextCancelled(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 1, Interval: 10 * time.Second})
	rl.Allow() // exhaust the single slot

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestRateLimiter_ConcurrentAllow(t *testing.T) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Interval: time.Second})
	allowed := make(chan bool, 200)
	for i := 0; i < 200; i++ {
		go func() { allowed <- rl.Allow() }()
	}
	count := 0
	for i := 0; i < 200; i++ {
		if <-allowed {
			count++
		}
	}
	if count > 100 {
		t.Errorf("expected at most 100 allowed, got %d", count)
	}
}
