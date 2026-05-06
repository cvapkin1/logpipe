package pipeline

import (
	"errors"
	"sync"
)

// ErrBufferFull is returned when the ring buffer has no capacity.
var ErrBufferFull = errors.New("pipeline: buffer is full")

// Buffer is a thread-safe fixed-capacity ring buffer for log lines.
type Buffer struct {
	mu       sync.Mutex
	data     []string
	head     int
	tail     int
	size     int
	capacity int
}

// NewBuffer creates a new ring buffer with the given capacity.
func NewBuffer(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = 256
	}
	return &Buffer{
		data:     make([]string, capacity),
		capacity: capacity,
	}
}

// Push adds a line to the buffer. Returns ErrBufferFull if at capacity.
func (b *Buffer) Push(line string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == b.capacity {
		return ErrBufferFull
	}
	b.data[b.tail] = line
	b.tail = (b.tail + 1) % b.capacity
	b.size++
	return nil
}

// Pop removes and returns the oldest line. Returns "", false if empty.
func (b *Buffer) Pop() (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == 0 {
		return "", false
	}
	line := b.data[b.head]
	b.head = (b.head + 1) % b.capacity
	b.size--
	return line, true
}

// Len returns the current number of items in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.size
}

// Cap returns the maximum capacity of the buffer.
func (b *Buffer) Cap() int {
	return b.capacity
}
