package pipeline

import (
	"fmt"
	"regexp"
)

// RedactConfig holds configuration for a single redaction rule.
type RedactConfig struct {
	// Pattern is the regular expression to match sensitive data.
	Pattern string
	// Placeholder replaces each match; defaults to "[REDACTED]".
	Placeholder string
}

// DefaultRedactPlaceholder is used when RedactConfig.Placeholder is empty.
const DefaultRedactPlaceholder = "[REDACTED]"

// NewRedactTransform compiles a RedactConfig into a TransformFunc.
// It returns an error if the pattern is invalid.
func NewRedactTransform(cfg RedactConfig) (TransformFunc, error) {
	if cfg.Pattern == "" {
		return nil, fmt.Errorf("redact: pattern must not be empty")
	}
	re, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return nil, fmt.Errorf("redact: invalid pattern %q: %w", cfg.Pattern, err)
	}
	placeholder := cfg.Placeholder
	if placeholder == "" {
		placeholder = DefaultRedactPlaceholder
	}
	return func(line string) string {
		return re.ReplaceAllString(line, placeholder)
	}, nil
}

// MustRedactTransform is like NewRedactTransform but panics on error.
// Intended for use in package-level variable declarations with known-good patterns.
func MustRedactTransform(cfg RedactConfig) TransformFunc {
	f, err := NewRedactTransform(cfg)
	if err != nil {
		panic(err)
	}
	return f
}
