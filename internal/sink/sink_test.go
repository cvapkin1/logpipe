package sink

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestStdoutSink_Write(t *testing.T) {
	var buf bytes.Buffer
	s := &StdoutSink{name: "test-stdout", w: &buf}

	if err := s.Write([]byte("hello log")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello log\n" {
		t.Errorf("expected 'hello log\\n', got %q", got)
	}
}

func TestStdoutSink_Name(t *testing.T) {
	s := NewStdoutSink("my-sink")
	if s.Name() != "my-sink" {
		t.Errorf("expected 'my-sink', got %q", s.Name())
	}
}

func TestStdoutSink_Close(t *testing.T) {
	s := NewStdoutSink("close-test")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected error on Close: %v", err)
	}
}

func TestTCPSink_WriteAndClose(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	received := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		received <- string(buf[:n])
	}()

	s, err := NewTCPSink("tcp-test", ln.Addr().String())
	if err != nil {
		t.Fatalf("NewTCPSink: %v", err)
	}
	defer s.Close()

	if err := s.Write([]byte("tcp log line")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got := <-received
	if got != "tcp log line\n" {
		t.Errorf("expected 'tcp log line\\n', got %q", got)
	}
}

func TestTCPSink_Name(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	s, err := NewTCPSink("tcp-named", ln.Addr().String())
	if err != nil {
		t.Fatalf("NewTCPSink: %v", err)
	}
	defer s.Close()

	if got := s.Name(); got != "tcp-named" {
		t.Errorf("expected 'tcp-named', got %q", got)
	}
}

func TestNewTCPSink_DialFailure(t *testing.T) {
	_, err := NewTCPSink("bad", "127.0.0.1:1")
	if err == nil {
		t.Error("expected error for unreachable address, got nil")
	}
	fmt.Println("expected error:", err)
}
