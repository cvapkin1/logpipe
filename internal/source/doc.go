// Package source provides log source implementations for logpipe.
//
// A Source reads raw log lines from an input (stdin, TCP socket, etc.) and
// emits them as strings on a channel for downstream processing by the
// pipeline router.
//
// Available source types:
//
//   - StdinSource  – reads newline-delimited log lines from os.Stdin.
//   - TCPSource    – listens on a TCP address and reads lines from every
//     accepted connection.
//
// Sources are managed through a Registry which enforces unique names and
// provides bulk lifecycle operations (CloseAll).
//
// Example usage:
//
//	src, err := source.NewTCPSource("app", ":5170")
//	if err != nil { ... }
//	out := make(chan string, 256)
//	src.Start(ctx, out)
//	for line := range out {
//		// forward line to pipeline
//	}
package source
