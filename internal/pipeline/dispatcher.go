package pipeline

import (
	"context"
	"log"
	"sync"

	"github.com/yourorg/logpipe/internal/source"
)

// Dispatcher reads log lines from a set of sources and forwards each
// line to a Router for rule-based delivery to sinks.
type Dispatcher struct {
	router  *Router
	sources []source.Source
	wg      sync.WaitGroup
}

// NewDispatcher creates a Dispatcher that uses the given Router to
// forward lines received from any registered source.
func NewDispatcher(r *Router) *Dispatcher {
	return &Dispatcher{router: r}
}

// AddSource registers a source whose lines will be dispatched.
func (d *Dispatcher) AddSource(s source.Source) {
	d.sources = append(d.sources, s)
}

// Run starts all registered sources and begins routing their output.
// It blocks until the context is cancelled and all sources have stopped.
func (d *Dispatcher) Run(ctx context.Context) {
	for _, s := range d.sources {
		ch, err := s.Start(ctx)
		if err != nil {
			log.Printf("dispatcher: failed to start source %q: %v", s.Name(), err)
			continue
		}
		d.wg.Add(1)
		go func(name string, lines <-chan string) {
			defer d.wg.Done()
			for line := range lines {
				if err := d.router.Route(ctx, line); err != nil {
					log.Printf("dispatcher: route error from source %q: %v", name, err)
				}
			}
		}(s.Name(), ch)
	}
	d.wg.Wait()
}
