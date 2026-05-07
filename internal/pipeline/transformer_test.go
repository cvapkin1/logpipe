package pipeline

import (
	"strings"
	"testing"
)

func TestTransformer_Apply_Empty(t *testing.T) {
	tr := NewTransformer()
	out, keep := tr.Apply("hello")
	if !keep || out != "hello" {
		t.Fatalf("expected (hello, true), got (%q, %v)", out, keep)
	}
}

func TestTransformer_Apply_SingleFunc(t *testing.T) {
	tr := NewTransformer(PrefixTransform("[LOG] "))
	out, keep := tr.Apply("message")
	if !keep || out != "[LOG] message" {
		t.Fatalf("unexpected result: %q, %v", out, keep)
	}
}

func TestTransformer_Apply_ChainedFuncs(t *testing.T) {
	tr := NewTransformer(
		TrimSpaceTransform(),
		PrefixTransform(">> "),
	)
	out, keep := tr.Apply("  hello world  ")
	if !keep || out != ">> hello world" {
		t.Fatalf("unexpected result: %q, %v", out, keep)
	}
}

func TestTransformer_Apply_DropOnEmpty(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	out, keep := tr.Apply("   ")
	if keep || out != "" {
		t.Fatalf("expected line to be dropped, got (%q, %v)", out, keep)
	}
}

func TestTransformer_Apply_StopsOnDrop(t *testing.T) {
	called := false
	tr := NewTransformer(
		TrimSpaceTransform(),
		func(line string) (string, bool) {
			called = true
			return line, true
		},
	)
	_, keep := tr.Apply("  ")
	if keep {
		t.Fatal("expected line to be dropped")
	}
	if called {
		t.Fatal("subsequent transform should not be called after drop")
	}
}

func TestMaxLengthTransform(t *testing.T) {
	tr := NewTransformer(MaxLengthTransform(5))
	out, keep := tr.Apply("hello world")
	if !keep || out != "hello" {
		t.Fatalf("expected truncation to 5 chars, got %q", out)
	}
}

func TestMaxLengthTransform_NoTruncation(t *testing.T) {
	tr := NewTransformer(MaxLengthTransform(100))
	out, keep := tr.Apply("short")
	if !keep || out != "short" {
		t.Fatalf("unexpected result: %q, %v", out, keep)
	}
}

func TestTimestampTransform_Prepends(t *testing.T) {
	tr := NewTransformer(TimestampTransform(""))
	out, keep := tr.Apply("event")
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if !strings.HasSuffix(out, " event") {
		t.Fatalf("expected timestamp prefix, got: %q", out)
	}
}
