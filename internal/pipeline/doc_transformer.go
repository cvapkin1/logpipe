// Package pipeline provides core log processing primitives for logpipe.
//
// # Transformer
//
// A Transformer applies a chain of transform functions to each log line
// before it is forwarded to a sink. Transform functions may modify, enrich,
// or drop a line entirely by returning an empty string.
//
// Built-in transforms:
//
//   - TrimSpaceTransform  – strips leading/trailing whitespace
//   - PrefixTransform     – prepends a fixed string to every line
//   - TimestampTransform  – prepends an RFC3339 UTC timestamp
//   - RedactTransform     – replaces matches of a regexp with a placeholder
//
// Example:
//
//	t := pipeline.NewTransformer(
//		pipeline.TrimSpaceTransform,
//		pipeline.PrefixTransform("[app] "),
//		pipeline.TimestampTransform,
//	)
//	out := t.Apply("  hello world  ") // "[app] 2024-01-01T00:00:00Z   hello world"
package pipeline
