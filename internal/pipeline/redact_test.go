package pipeline

import (
	"strings"
	"testing"
)

func TestNewRedactTransform_EmptyPattern(t *testing.T) {
	_, err := NewRedactTransform(RedactConfig{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNewRedactTransform_InvalidPattern(t *testing.T) {
	_, err := NewRedactTransform(RedactConfig{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNewRedactTransform_DefaultPlaceholder(t *testing.T) {
	fn, err := NewRedactTransform(RedactConfig{Pattern: `\d{4}-\d{4}-\d{4}-\d{4}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fn("card: 1234-5678-9012-3456 processed")
	if strings.Contains(out, "1234-5678-9012-3456") {
		t.Errorf("sensitive data not redacted: %q", out)
	}
	if !strings.Contains(out, DefaultRedactPlaceholder) {
		t.Errorf("expected placeholder %q in output: %q", DefaultRedactPlaceholder, out)
	}
}

func TestNewRedactTransform_CustomPlaceholder(t *testing.T) {
	fn, err := NewRedactTransform(RedactConfig{
		Pattern:     `password=\S+`,
		Placeholder: "password=***",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := fn("login password=secret123 ok")
	if strings.Contains(out, "secret123") {
		t.Errorf("secret not redacted: %q", out)
	}
	if !strings.Contains(out, "password=***") {
		t.Errorf("expected custom placeholder in output: %q", out)
	}
}

func TestNewRedactTransform_NoMatch(t *testing.T) {
	fn, err := NewRedactTransform(RedactConfig{Pattern: `token=\S+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "ordinary log line without sensitive data"
	out := fn(input)
	if out != input {
		t.Errorf("expected unchanged output, got %q", out)
	}
}

func TestMustRedactTransform_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid pattern, got none")
		}
	}()
	MustRedactTransform(RedactConfig{Pattern: "[bad"})
}

func TestRedactTransform_IntegrationWithTransformer(t *testing.T) {
	redact, err := NewRedactTransform(RedactConfig{Pattern: `\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr := NewTransformer(TrimSpaceTransform, redact)
	out := tr.Apply("  user user@example.com logged in  ")
	if strings.Contains(out, "user@example.com") {
		t.Errorf("email not redacted in transformer chain: %q", out)
	}
	if !strings.Contains(out, DefaultRedactPlaceholder) {
		t.Errorf("placeholder missing in transformer chain output: %q", out)
	}
}
