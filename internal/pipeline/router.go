// Package pipeline handles log routing between sources and destinations.
package pipeline

import (
	"fmt"
	"sync"

	"github.com/yourorg/logpipe/internal/filter"
)

// LogEntry represents a single log line with metadata.
type LogEntry struct {
	Source  string
	Level   filter.Level
	Message string
	Raw     []byte
}

// Destination is a sink that receives matched log entries.
type Destination interface {
	Write(entry LogEntry) error
	Name() string
}

// Route binds a set of filter rules to a destination.
type Route struct {
	Rules       []*filter.Rule
	Destination Destination
}

// Router dispatches log entries to destinations based on matching routes.
type Router struct {
	mu     sync.RWMutex
	routes []*Route
}

// NewRouter creates an empty Router.
func NewRouter() *Router {
	return &Router{}
}

// AddRoute registers a route with the router.
func (r *Router) AddRoute(route *Route) error {
	if route == nil {
		return fmt.Errorf("route must not be nil")
	}
	if route.Destination == nil {
		return fmt.Errorf("route destination must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
	return nil
}

// Dispatch sends an entry to every destination whose route matches it.
// Returns the number of destinations the entry was forwarded to.
func (r *Router) Dispatch(entry LogEntry) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	forwarded := 0
	var firstErr error

	for _, route := range r.routes {
		if matchesAll(route.Rules, entry) {
			if err := route.Destination.Write(entry); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("destination %q: %w", route.Destination.Name(), err)
			}
			forwarded++
		}
	}
	return forwarded, firstErr
}

// matchesAll returns true when the entry satisfies every rule in the slice.
// An empty rule slice matches all entries.
func matchesAll(rules []*filter.Rule, entry LogEntry) bool {
	for _, rule := range rules {
		if !rule.Match(entry.Level, entry.Message) {
			return false
		}
	}
	return true
}
