package pipeline

import (
	"testing"
)

func TestNewSampler_DefaultConfig(t *testing.T) {
	cfg := DefaultSamplerConfig()
	if cfg.Rate != 1.0 {
		t.Fatalf("expected default rate 1.0, got %v", cfg.Rate)
	}
}

func TestNewSampler_InvalidRate_Zero(t *testing.T) {
	_, err := NewSampler(SamplerConfig{Rate: 0})
	if err == nil {
		t.Fatal("expected error for rate=0, got nil")
	}
}

func TestNewSampler_InvalidRate_Negative(t *testing.T) {
	_, err := NewSampler(SamplerConfig{Rate: -0.5})
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNewSampler_InvalidRate_OverOne(t *testing.T) {
	_, err := NewSampler(SamplerConfig{Rate: 1.1})
	if err == nil {
		t.Fatal("expected error for rate > 1, got nil")
	}
}

func TestSampler_Allow_RateOne_AllPass(t *testing.T) {
	s, err := NewSampler(SamplerConfig{Rate: 1.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 20; i++ {
		if !s.Allow() {
			t.Fatalf("expected Allow()=true at call %d with rate=1.0", i+1)
		}
	}
}

func TestSampler_Allow_RateHalf_KeepsEveryOther(t *testing.T) {
	s, err := NewSampler(SamplerConfig{Rate: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < 10; i++ {
		if s.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed out of 10 at rate=0.5, got %d", allowed)
	}
}

func TestSampler_Allow_RateTenth_KeepsOneInTen(t *testing.T) {
	s, err := NewSampler(SamplerConfig{Rate: 0.1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < 100; i++ {
		if s.Allow() {
			allowed++
		}
	}
	if allowed != 10 {
		t.Fatalf("expected 10 allowed out of 100 at rate=0.1, got %d", allowed)
	}
}

func TestSampler_Reset_RestartsCounter(t *testing.T) {
	s, err := NewSampler(SamplerConfig{Rate: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// consume first slot
	s.Allow()
	s.Allow()
	s.Reset()
	// after reset the first call should be allowed again
	if !s.Allow() {
		t.Fatal("expected Allow()=true immediately after Reset()")
	}
}
