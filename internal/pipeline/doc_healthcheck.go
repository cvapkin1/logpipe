// Package pipeline provides log pipeline primitives for logpipe.
//
// # Health Checking
//
// HealthChecker evaluates a live Metrics snapshot to determine whether the
// pipeline is operating normally. Three statuses are defined:
//
//   - StatusHealthy   – drop rate is below the configured threshold and
//     messages have been received within the stale window.
//
//   - StatusDegraded  – the fraction of dropped messages meets or exceeds
//     the configured degradedDropRate (default 0.10). This typically signals
//     backpressure or a slow downstream sink.
//
//   - StatusUnhealthy – no message has been received within staleDuration
//     (default 30 s). This may indicate a broken source or network partition.
//
// Usage:
//
//	m := pipeline.NewMetrics()
//	hc := pipeline.NewHealthChecker(m, 0.05, 15*time.Second)
//	report := hc.Check()
//	if report.Status != pipeline.StatusHealthy {
//		log.Printf("pipeline %s: %s", report.Status, report.Message)
//	}
package pipeline
