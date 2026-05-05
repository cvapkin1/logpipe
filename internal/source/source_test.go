package source

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestStdinSource_Name(t *testing.T) {
	s := NewStdinSource("stdin-test")
	if s.Name() != "stdin-test" {
		t.Fatalf("expected 'stdin-test', got %q", s.Name())
	}
}

func TestStdinSource_Start(t *testing.T) {
	s := &StdinSource{name: "test", r: strings.NewReader("hello\nworld\n")}
	out := make(chan string, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Start(ctx, out); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	var lines []string
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
loop:
	for {
		select {
		case line := <-out:
			lines = append(lines, line)
			if len(lines) == 2 {
				break loop
			}
		case <-timer.C:
			break loop
		}
	}

	if len(lines) != 2 || lines[0] != "hello" || lines[1] != "world" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestTCPSource_NameAndClose(t *testing.T) {
	s, err := NewTCPSource("tcp-test", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("NewTCPSource: %v", err)
	}
	if s.Name() != "tcp-test" {
		t.Fatalf("expected 'tcp-test', got %q", s.Name())
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestTCPSource_ReceivesLines(t *testing.T) {
	s, err := NewTCPSource("tcp-recv", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("NewTCPSource: %v", err)
	}
	defer s.Close()

	addr := s.listener.Addr().String()
	out := make(chan string, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Start(ctx, out); err != nil {
		t.Fatalf("Start: %v", err)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	fmt.Fprintln(conn, "log line one")
	fmt.Fprintln(conn, "log line two")
	conn.Close()

	var got []string
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
loop:
	for {
		select {
		case line := <-out:
			got = append(got, line)
			if len(got) == 2 {
				break loop
			}
		case <-timer.C:
			break loop
		}
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}

	_ = bufio.NewScanner(strings.NewReader("")) // suppress unused import
}

func TestNewTCPSource_BadAddress(t *testing.T) {
	_, err := NewTCPSource("bad", "invalid-address")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
