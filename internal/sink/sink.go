// Package sink provides destination writers for log forwarding.
package sink

import (
	"fmt"
	"io"
	"net"
	"os"
)

// Sink is a destination that can receive log lines.
type Sink interface {
	Write(line []byte) error
	Close() error
	Name() string
}

// StdoutSink writes log lines to standard output.
type StdoutSink struct {
	name string
	w    io.Writer
}

// NewStdoutSink creates a new StdoutSink with the given name.
func NewStdoutSink(name string) *StdoutSink {
	return &StdoutSink{name: name, w: os.Stdout}
}

func (s *StdoutSink) Write(line []byte) error {
	_, err := fmt.Fprintf(s.w, "%s\n", line)
	return err
}

func (s *StdoutSink) Close() error { return nil }
func (s *StdoutSink) Name() string { return s.name }

// TCPSink forwards log lines to a remote TCP endpoint.
type TCPSink struct {
	name    string
	address string
	conn    net.Conn
}

// NewTCPSink creates a TCPSink and establishes a connection to addr.
func NewTCPSink(name, addr string) (*TCPSink, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("sink %q: dial %s: %w", name, addr, err)
	}
	return &TCPSink{name: name, address: addr, conn: conn}, nil
}

func (s *TCPSink) Write(line []byte) error {
	_, err := fmt.Fprintf(s.conn, "%s\n", line)
	return err
}

func (s *TCPSink) Close() error { return s.conn.Close() }
func (s *TCPSink) Name() string { return s.name }
