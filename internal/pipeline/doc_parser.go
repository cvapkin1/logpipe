// Package pipeline provides stream-processing primitives for logpipe.
//
// # Parser
//
// Parser extracts structured fields from raw log lines using a compiled
// regular expression with named capture groups.
//
// Recognised special group names:
//
//   - level     – log severity, uppercased and stored in ParsedLine.Level
//   - message / msg – human-readable log body stored in ParsedLine.Message
//   - ts / timestamp / time – RFC 3339 timestamp stored in ParsedLine.Timestamp
//
// Any other named groups are collected into ParsedLine.Fields.
//
// Example:
//
//	pattern := `(?P<ts>\S+) (?P<level>\w+) (?P<message>.+)`
//	p, err := pipeline.NewParser(pattern)
//	if err != nil {
//		log.Fatal(err)
//	}
//	pl := p.Parse("2024-01-01T00:00:00Z INFO server started")
//	fmt.Println(pl.Level)   // INFO
//	fmt.Println(pl.Message) // server started
package pipeline
