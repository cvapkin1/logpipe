package pipeline

import (
	"errors"
	"sync/atomic"
)

// SamplerConfig holds configuration for the log sampler.
type SamplerConfig struct {
	// Rate is the fraction of messages to keep, in the range (0, 1].
	// A value of 1.0 keeps all messages; 0.1 keeps roughly 1 in 10.
	Rate float64
}

// DefaultSamplerConfig returns a SamplerConfig that keeps every message.
func DefaultSamplerConfig() SamplerConfig {
	return SamplerConfig{Rate: 1.0}
}

// Sampler deterministically samples log lines using a counter-based
// approach so that exactly 1-in-N messages are forwarded, where N is
// derived from the configured rate.
//
// It is safe for concurrent use.
type Sampler struct {
	every uint64 // keep 1 out of every N messages
	counter atomic.Uint64
}

// NewSampler creates a Sampler from the given config.
// Rate is clamped to (0, 1]; values outside that range return an error.
func NewSampler(cfg SamplerConfig) (*Sampler, error) {
	if cfg.Rate <= 0 || cfg.Rate > 1.0 {
		return nil, errors.New("sampler: rate must be in the range (0, 1]")
	}
	every := uint64(1.0 / cfg.Rate)
	if every == 0 {
		every = 1
	}
	return &Sampler{every: every}, nil
}

// Allow returns true if the current message should be forwarded.
// It increments an internal counter on every call.
func (s *Sampler) Allow() bool {
	n := s.counter.Add(1)
	return n%s.every == 1
}

// Reset resets the internal counter to zero.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}
