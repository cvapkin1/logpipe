package filter

import (
	"regexp"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelUnknown
)

// Rule defines filtering criteria for log lines.
type Rule struct {
	MinLevel   Level
	Contains   string
	Pattern    *regexp.Regexp
}

// ParseLevel converts a string level name to a Level constant.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelUnknown
	}
}

// NewRule constructs a Rule from a minimum level string, a substring, and an
// optional regex pattern string. Returns an error if the pattern is invalid.
func NewRule(minLevel, contains, pattern string) (*Rule, error) {
	r := &Rule{
		MinLevel: ParseLevel(minLevel),
		Contains: contains,
	}
	if pattern != "" {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		r.Pattern = compiled
	}
	return r, nil
}

// Match reports whether the given log line satisfies the rule.
// lineLevel is the parsed level of the incoming log line.
func (r *Rule) Match(line string, lineLevel Level) bool {
	if lineLevel != LevelUnknown && lineLevel < r.MinLevel {
		return false
	}
	if r.Contains != "" && !strings.Contains(line, r.Contains) {
		return false
	}
	if r.Pattern != nil && !r.Pattern.MatchString(line) {
		return false
	}
	return true
}
