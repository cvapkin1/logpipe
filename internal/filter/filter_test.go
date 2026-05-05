package filter

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARNING", LevelWarn},
		{"error", LevelError},
		{"unknown", LevelUnknown},
		{"", LevelUnknown},
	}
	for _, c := range cases {
		got := ParseLevel(c.input)
		if got != c.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestNewRule_InvalidPattern(t *testing.T) {
	_, err := NewRule("info", "", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex pattern, got nil")
	}
}

func TestRule_Match_MinLevel(t *testing.T) {
	r, err := NewRule("warn", "", "")
	if err != nil {
		t.Fatalf("NewRule: %v", err)
	}
	if r.Match("some debug line", LevelDebug) {
		t.Error("expected debug line to NOT match warn rule")
	}
	if !r.Match("some error line", LevelError) {
		t.Error("expected error line to match warn rule")
	}
}

func TestRule_Match_Contains(t *testing.T) {
	r, err := NewRule("", "timeout", "")
	if err != nil {
		t.Fatalf("NewRule: %v", err)
	}
	if !r.Match("connection timeout reached", LevelUnknown) {
		t.Error("expected line containing 'timeout' to match")
	}
	if r.Match("all systems operational", LevelUnknown) {
		t.Error("expected line without 'timeout' to NOT match")
	}
}

func TestRule_Match_Pattern(t *testing.T) {
	r, err := NewRule("", "", `error\s+code:\s*\d+`)
	if err != nil {
		t.Fatalf("NewRule: %v", err)
	}
	if !r.Match("error code: 503", LevelUnknown) {
		t.Error("expected pattern match")
	}
	if r.Match("everything is fine", LevelUnknown) {
		t.Error("expected no pattern match")
	}
}

func TestRule_Match_Combined(t *testing.T) {
	r, err := NewRule("error", "disk", `disk\s+full`)
	if err != nil {
		t.Fatalf("NewRule: %v", err)
	}
	// Passes all criteria
	if !r.Match("disk full error", LevelError) {
		t.Error("expected combined rule to match")
	}
	// Fails level
	if r.Match("disk full", LevelInfo) {
		t.Error("expected combined rule to fail on level")
	}
}
