package source

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

// Source reads log lines and sends them to an output channel.
type Source interface {
	Name() string
	Start(ctx context.Context, out chan<- string) error
	Close() error
}

// StdinSource reads log lines from standard input.
type StdinSource struct {
	name string
	r    io.Reader
}

// NewStdinSource creates a Source that reads from stdin.
func NewStdinSource(name string) *StdinSource {
	return &StdinSource{name: name, r: os.Stdin}
}

func (s *StdinSource) Name() string { return s.name }

func (s *StdinSource) Start(ctx context.Context, out chan<- string) error {
	go func() {
		scanner := bufio.NewScanner(s.r)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case out <- scanner.Text():
			}
		}
	}()
	return nil
}

func (s *StdinSource) Close() error { return nil }

// TCPSource listens on a TCP address and reads log lines from connections.
type TCPSource struct {
	name     string
	address  string
	listener net.Listener
}

// NewTCPSource creates a Source that listens for TCP connections.
func NewTCPSource(name, address string) (*TCPSource, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("source %s: listen %s: %w", name, address, err)
	}
	return &TCPSource{name: name, address: address, listener: ln}, nil
}

func (s *TCPSource) Name() string { return s.name }

func (s *TCPSource) Start(ctx context.Context, out chan<- string) error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return
			}
			go s.handleConn(ctx, conn, out)
		}
	}()
	return nil
}

func (s *TCPSource) handleConn(ctx context.Context, conn net.Conn, out chan<- string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case out <- scanner.Text():
		}
	}
}

func (s *TCPSource) Close() error { return s.listener.Close() }
