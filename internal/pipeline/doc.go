// Package pipeline provides log routing and forwarding infrastructure
// for logpipe. It connects sources to sinks via configurable rules,
// allowing log lines to be filtered and forwarded to one or more
// destinations based on matching criteria.
//
// The central type is Router, which maintains a set of routes. Each
// route pairs a filter Rule with a Destination sink. Incoming log lines
// are evaluated against every route; matching routes forward the line
// to their associated sink.
//
// Usage:
//
//	router := pipeline.NewRouter()
//	router.AddRoute(rule, sink)
//	router.Route(ctx, line)
package pipeline
