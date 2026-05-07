package pipeline

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Checkpoint records the last successfully processed offset or cursor
// for a named source, enabling resume-after-restart semantics.
type Checkpoint struct {
	mu      sync.Mutex
	path    string
	entries map[string]CheckpointEntry
}

// CheckpointEntry holds the persisted state for a single source.
type CheckpointEntry struct {
	Source    string    `json:"source"`
	Offset    int64     `json:"offset"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCheckpoint creates a Checkpoint backed by the given file path.
// If the file exists its contents are loaded; a missing file is not an error.
func NewCheckpoint(path string) (*Checkpoint, error) {
	cp := &Checkpoint{
		path:    path,
		entries: make(map[string]CheckpointEntry),
	}
	if err := cp.load(); err != nil {
		return nil, err
	}
	return cp, nil
}

// Set records the latest offset for the named source and flushes to disk.
func (c *Checkpoint) Set(source string, offset int64) error {
	c.mu.Lock()
	c.entries[source] = CheckpointEntry{
		Source:    source,
		Offset:    offset,
		UpdatedAt: time.Now().UTC(),
	}
	c.mu.Unlock()
	return c.flush()
}

// Get returns the last saved offset for source, or 0 if not found.
func (c *Checkpoint) Get(source string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.entries[source]; ok {
		return e.Offset
	}
	return 0
}

// Delete removes the checkpoint entry for source and flushes to disk.
func (c *Checkpoint) Delete(source string) error {
	c.mu.Lock()
	delete(c.entries, source)
	c.mu.Unlock()
	return c.flush()
}

func (c *Checkpoint) load() error {
	data, err := os.ReadFile(c.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Unmarshal(data, &c.entries)
}

func (c *Checkpoint) flush() error {
	c.mu.Lock()
	data, err := json.MarshalIndent(c.entries, "", "  ")
	c.mu.Unlock()
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}
