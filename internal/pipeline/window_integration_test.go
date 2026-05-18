package pipeline

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// TestWindow_Integration_ConcurrentAdds verifies that concurrent writers do
// not corrupt internal state and that all lines are eventually flushed.
func TestWindow_Integration_ConcurrentAdds(t *testing.T) {
	w, err := NewWindow(WindowConfig{Size: 50 * time.Millisecond, MaxLines: 20})
	if err != nil {
		t.Fatalf("NewWindow: %v", err)
	}

	var (
		mu      sync.Mutex
		collect []string
		wg      sync.WaitGroup
	)

	const writers = 5
	const linesEach = 10

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < linesEach; j++ {
				line := strings.Repeat("x", id+j+1)
				if batch, ok := w.Add(line); ok {
					mu.Lock()
					collect = append(collect, batch...)
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()

	// Flush remaining.
	if tail := w.Close(); len(tail) > 0 {
		mu.Lock()
		collect = append(collect, tail...)
		mu.Unlock()
	}

	total := writers * linesEach
	if len(collect) != total {
		t.Fatalf("expected %d total lines, got %d", total, len(collect))
	}
}

// TestWindow_Integration_WithDispatcher verifies that a Window can feed
// batched lines into a Dispatcher without data loss.
func TestWindow_Integration_WithDispatcher(t *testing.T) {
	w, _ := NewWindow(WindowConfig{Size: 5 * time.Second, MaxLines: 3})

	var received []string
	var mu sync.Mutex

	d := NewDispatcher()
	d.Subscribe(func(line string) {
		mu.Lock()
		received = append(received, line)
		mu.Unlock()
	})

	flushBatch := func(batch []string) {
		for _, l := range batch {
			d.Dispatch(l)
		}
	}

	input := []string{"alpha", "beta", "gamma", "delta"}
	for _, line := range input {
		if batch, ok := w.Add(line); ok {
			flushBatch(batch)
		}
	}
	flushBatch(w.Close())

	mu.Lock()
	defer mu.Unlock()
	if len(received) != len(input) {
		t.Fatalf("expected %d dispatched lines, got %d", len(input), len(received))
	}
}
