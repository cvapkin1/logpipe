// Package sink defines the Sink interface and concrete implementations
// for forwarding log lines to various destinations.
//
// # Available Sinks
//
// StdoutSink writes each log line to standard output, suitable for
// local development and container log collection via stdout.
//
// TCPSink dials a remote TCP address and streams log lines over the
// connection. Useful for forwarding to log aggregators such as Logstash
// or a custom TCP listener.
//
// # Registry
//
// Registry provides a thread-safe store for named Sink instances.
// Use Register to add sinks, Get to look them up by name, and
// CloseAll to gracefully shut down all connections on exit.
//
// Example:
//
//	reg := sink.NewRegistry()
//	reg.Register(sink.NewStdoutSink("console"))
//	s, ok := reg.Get("console")
//	if ok {
//		s.Write([]byte("log message"))
//	}
package sink
