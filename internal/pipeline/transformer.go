package pipeline

import (
	"strings"
	"time"
)

// TransformFunc is a function that transforms a log line.
// It returns the transformed line and a boolean indicating whether the line should be kept.
type TransformFunc func(line string) (string, bool)

// Transformer applies a chain of transformation functions to log lines.
type Transformer struct {
	funcs []TransformFunc
}

// NewTransformer creates a new Transformer with the given transform functions.
func NewTransformer(fns ...TransformFunc) *Transformer {
	return &Transformer{funcs: fns}
}

// Apply runs the line through all transform functions in order.
// If any function signals the line should be dropped (returns false), Apply returns ("", false).
func (t *Transformer) Apply(line string) (string, bool) {
	current := line
	for _, fn := range t.funcs {
		result, keep := fn(current)
		if !keep {
			return "", false
		}
		current = result
	}
	return current, true
}

// TrimSpaceTransform returns a TransformFunc that trims leading and trailing whitespace.
func TrimSpaceTransform() TransformFunc {
	return func(line string) (string, bool) {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			return "", false
		}
		return trimmed, true
	}
}

// PrefixTransform returns a TransformFunc that prepends a static prefix to each line.
func PrefixTransform(prefix string) TransformFunc {
	return func(line string) (string, bool) {
		return prefix + line, true
	}
}

// TimestampTransform returns a TransformFunc that prepends the current UTC timestamp.
func TimestampTransform(layout string) TransformFunc {
	if layout == "" {
		layout = time.RFC3339
	}
	return func(line string) (string, bool) {
		return time.Now().UTC().Format(layout) + " " + line, true
	}
}

// MaxLengthTransform returns a TransformFunc that truncates lines exceeding maxLen bytes.
func MaxLengthTransform(maxLen int) TransformFunc {
	return func(line string) (string, bool) {
		if maxLen > 0 && len(line) > maxLen {
			return line[:maxLen], true
		}
		return line, true
	}
}
