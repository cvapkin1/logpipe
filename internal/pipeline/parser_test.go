package pipeline

import (
	"testing"
	"time"
)

const samplePattern = `(?P<ts>\S+) (?P<level>\w+) (?P<message>.+)`

func TestNewParser_EmptyPattern(t *testing.T) {
	_, err := NewParser("")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewParser_InvalidRegex(t *testing.T) {
	_, err := NewParser("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNewParser_NoNamedGroups(t *testing.T) {
	_, err := NewParser(`(\w+) (\w+)`)
	if err == nil {
		t.Fatal("expected error when no named groups present")
	}
}

func TestNewParser_Valid(t *testing.T) {
	p, err := NewParser(samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestParser_Parse_MatchesLevel(t *testing.T) {
	p, _ := NewParser(samplePattern)
	pl := p.Parse("2024-01-01T00:00:00Z INFO server started")
	if pl.Level != "INFO" {
		t.Errorf("expected level INFO, got %q", pl.Level)
	}
}

func TestParser_Parse_MatchesMessage(t *testing.T) {
	p, _ := NewParser(samplePattern)
	pl := p.Parse("2024-01-01T00:00:00Z ERROR disk full")
	if pl.Message != "disk full" {
		t.Errorf("expected message 'disk full', got %q", pl.Message)
	}
}

func TestParser_Parse_MatchesTimestamp(t *testing.T) {
	p, _ := NewParser(samplePattern)
	pl := p.Parse("2024-06-15T12:00:00Z WARN high memory")
	expected := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	if !pl.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, pl.Timestamp)
	}
}

func TestParser_Parse_NoMatch_ReturnsRaw(t *testing.T) {
	p, _ := NewParser(samplePattern)
	raw := "not a structured line"
	pl := p.Parse(raw)
	if pl.Raw != raw {
		t.Errorf("expected Raw=%q, got %q", raw, pl.Raw)
	}
	if pl.Level != "" {
		t.Errorf("expected empty Level on no-match, got %q", pl.Level)
	}
}

func TestParser_Parse_ExtraFields(t *testing.T) {
	pattern := `(?P<level>\w+) (?P<message>.+) host=(?P<host>\S+)`
	p, err := NewParser(pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pl := p.Parse("INFO service ready host=web-01")
	if pl.Fields["host"] != "web-01" {
		t.Errorf("expected host=web-01, got %q", pl.Fields["host"])
	}
}

func TestParser_Parse_RawAlwaysSet(t *testing.T) {
	p, _ := NewParser(samplePattern)
	line := "2024-01-01T00:00:00Z DEBUG connecting"
	pl := p.Parse(line)
	if pl.Raw != line {
		t.Errorf("expected Raw=%q, got %q", line, pl.Raw)
	}
}
