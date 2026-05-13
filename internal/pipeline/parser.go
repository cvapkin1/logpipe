package pipeline

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ParsedLine holds structured fields extracted from a raw log line.
type ParsedLine struct {
	Raw       string
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
}

// Parser extracts structured fields from raw log lines using a named-group regex.
type Parser struct {
	pattern *regexp.Regexp
	fields  []string
}

// NewParser compiles the given regex pattern. The pattern must use named capture
// groups to identify fields (e.g. (?P<level>INFO)). Returns an error if the
// pattern is invalid or contains no named groups.
func NewParser(pattern string) (*Parser, error) {
	if strings.TrimSpace(pattern) == "" {
		return nil, fmt.Errorf("parser: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("parser: invalid pattern: %w", err)
	}
	names := re.SubexpNames()
	var named []string
	for _, n := range names {
		if n != "" {
			named = append(named, n)
		}
	}
	if len(named) == 0 {
		return nil, fmt.Errorf("parser: pattern must contain at least one named capture group")
	}
	return &Parser{pattern: re, fields: named}, nil
}

// Parse attempts to match line against the compiled pattern and returns a
// ParsedLine. If the line does not match, Raw is populated and all other
// fields are left at their zero values.
func (p *Parser) Parse(line string) ParsedLine {
	result := ParsedLine{Raw: line, Fields: make(map[string]string)}
	match := p.pattern.FindStringSubmatch(line)
	if match == nil {
		return result
	}
	for i, name := range p.pattern.SubexpNames() {
		if name == "" || i >= len(match) {
			continue
		}
		val := match[i]
		switch strings.ToLower(name) {
		case "level":
			result.Level = strings.ToUpper(val)
		case "message", "msg":
			result.Message = val
		case "ts", "timestamp", "time":
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				result.Timestamp = t
			}
		default:
			result.Fields[name] = val
		}
	}
	if result.Message == "" {
		result.Message = line
	}
	return result
}
