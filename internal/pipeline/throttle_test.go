package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestNewThrottle_Defaults(t *testing.T) {
	th := NewThrottle(ThrottleConfig{})
	if th.cfg.Rate != 1000 {
		t.Fatalf("expected default rate 1000, got %d", th.cfg.Rate)
	}
	if th.cfg.Window != time.Second {
		t.Fatalf("expected default window 1s, got %v", th.cfg.Window)
	}
}

func TestThrottle_AllowsUpToRate(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 5, Window: time.Second, BurstSize: 0})
	for i := 0; i < 5; i++ {
		if !th.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if th.Allow() {
		t.Fatal("expected Allow()=false after rate exceeded")
	}
}

func TestThrottle_ResetsAfterWindow(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 2, Window: 20 * time.Millisecond, BurstSize: 0})
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected drop before window reset")
	}
	time.Sleep(25 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected Allow()=true after window reset")
	}
}

func TestThrottle_BurstAllowsExtra(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 3, Window: time.Second, BurstSize: 2})
	allowed := 0
	for i := 0; i < 10; i++ {
		if th.Allow() {
			allowed++
		}
	}
	// Should allow Rate + BurstSize = 5
	if allowed != 5 {
		t.Fatalf("expected 5 allowed with burst, got %d", allowed)
	}
}

func TestThrottle_Reset_ClearsState(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 2, Window: time.Second, BurstSize: 0})
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected drop before reset")
	}
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected Allow()=true after Reset")
	}
}

func TestThrottle_AllowCtx_CancelledReturnsFalse(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 100, Window: time.Second, BurstSize: 0})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if th.AllowCtx(ctx) {
		t.Fatal("expected false when context already cancelled")
	}
}

func TestThrottle_AllowCtx_ActiveContext(t *testing.T) {
	th := NewThrottle(ThrottleConfig{Rate: 10, Window: time.Second, BurstSize: 0})
	ctx := context.Background()
	if !th.AllowCtx(ctx) {
		t.Fatal("expected true with active context and room in rate")
	}
}

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.Rate <= 0 {
		t.Fatal("default Rate must be positive")
	}
	if cfg.Window <= 0 {
		t.Fatal("default Window must be positive")
	}
	if cfg.BurstSize < 0 {
		t.Fatal("default BurstSize must not be negative")
	}
}
