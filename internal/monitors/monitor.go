package monitors

import "github.com/JinFuuMugen/ya_go_metrics/internal/sender"

// Monitor defines a common interface for metric collectors.
type Monitor interface {
	// Collect gathers metrics from the source and stores them internally.
	Collect() error

	// Dump sends collected metrics using the configured Sender.
	Dump() error

	// SetProcessor sets the Sender used to send collected metrics.
	SetProcessor(p sender.Sender)
}

// RuntimeMonitor extends Monitor with runtime metric collection.
type RuntimeMonitor interface {
	Monitor

	// CollectRuntimeMetrics collects metrics from the Go runtime.
	CollectRuntimeMetrics()
}

// GopsutilMonitor extends Monitor with gopsutil-based metric collection.
type GopsutilMonitor interface {
	Monitor

	// CollectGopsutil collects metrics using the gopsutil library.
	CollectGopsutil() error
}
