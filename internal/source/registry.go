package source

import (
	"fmt"
	"sync"
)

// Registry holds named Sources.
type Registry struct {
	mu      sync.RWMutex
	sources map[string]Source
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{sources: make(map[string]Source)}
}

// Register adds a source under its name. Returns an error if the name is already taken.
func (r *Registry) Register(s Source) error {
	if s == nil {
		return fmt.Errorf("source: cannot register nil source")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.sources[s.Name()]; exists {
		return fmt.Errorf("source: name %q already registered", s.Name())
	}
	r.sources[s.Name()] = s
	return nil
}

// Get returns the Source registered under name, or an error if not found.
func (r *Registry) Get(name string) (Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sources[name]
	if !ok {
		return nil, fmt.Errorf("source: %q not found", name)
	}
	return s, nil
}

// Names returns a sorted slice of all registered source names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.sources))
	for n := range r.sources {
		names = append(names, n)
	}
	return names
}

// CloseAll closes every registered source, collecting errors.
func (r *Registry) CloseAll() []error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var errs []error
	for _, s := range r.sources {
		if err := s.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close source %q: %w", s.Name(), err))
		}
	}
	return errs
}
