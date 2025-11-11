package monitors

import "github.com/JinFuuMugen/ya_go_metrics/internal/sender"

type Monitor interface {
	Collect() error
	Dump() error
	SetProcessor(p sender.Sender)
}

type RuntimeMonitor interface {
	Monitor
	CollectRuntimeMetrics()
}

type GopsutilMonitor interface {
	Monitor
	CollectGopsutil() error
}
