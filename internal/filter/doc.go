// Package filter provides log line filtering primitives for logpipe.
//
// A Rule combines three optional criteria:
//
//   - MinLevel: only lines at or above this severity are forwarded.
//   - Contains: only lines containing this substring are forwarded.
//   - Pattern: only lines matching this regular expression are forwarded.
//
// All specified criteria must be satisfied for a line to match (logical AND).
// Criteria that are left at their zero value are ignored, making it easy to
// compose simple or complex rules from configuration.
//
// Example usage:
//
//	rule, err := filter.NewRule("warn", "timeout", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if rule.Match(line, filter.ParseLevel(detectedLevel)) {
//		// forward the line
//	}
package filter
