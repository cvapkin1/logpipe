package sink

import (
	"fmt"
	"sync"
)

// Registry holds named sinks and provides thread-safe access.
type Registry struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{sinks: make(map[string]Sink)}
}

// Register adds a sink to the registry. Returns an error if the name is taken.
func (r *Registry) Register(s Sink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.sinks[s.Name()]; exists {
		return fmt.Errorf("sink %q already registered", s.Name())
	}
	r.sinks[s.Name()] = s
	return nil
}

// Get retrieves a sink by name.
func (r *Registry) Get(name string) (Sink, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sinks[name]
	return s, ok
}

// CloseAll closes all registered sinks and returns the first error encountered.
func (r *Registry) CloseAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var firstErr error
	for _, s := range r.sinks {
		if err := s.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Names returns all registered sink names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.sinks))
	for n := range r.sinks {
		names = append(names, n)
	}
	return names
}
