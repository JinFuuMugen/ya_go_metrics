package monitors

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

//generate:reset
//go:generate go run ../../cmd/reset
type runtimeMonitor struct {
	Storage   storage.Storage
	Processor sender.Sender
}

// NewRuntimeMonitor creates a new runtime-based monitor.
func NewRuntimeMonitor(s storage.Storage, p sender.Sender) RuntimeMonitor {
	return &runtimeMonitor{Storage: s, Processor: p}
}

// Collect collects runtime metrics.
func (m *runtimeMonitor) Collect() error {
	m.CollectRuntimeMetrics()
	return nil
}

// CollectRuntimeMetrics collects metrics from the Go runtime.
func (m *runtimeMonitor) CollectRuntimeMetrics() {
	m.collectRuntime()
	m.collectRuntimeSystem()
}

// Dump sends collected metrics using the configured Sender.
func (m *runtimeMonitor) Dump() error {
	c := m.Storage.GetCounters()
	g := m.Storage.GetGauges()
	err := m.Processor.Process(c, g)
	if err != nil {
		return fmt.Errorf("error dumping metric: %w", err)
	}
	return nil
}

// SetProcessor updates the Sender used for dumping metrics.
func (m *runtimeMonitor) SetProcessor(p sender.Sender) {
	m.Processor = p
}

func (m *runtimeMonitor) collectRuntime() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Storage.SetGauge("BuckHashSys", float64(rtm.BuckHashSys))
	m.Storage.SetGauge("Alloc", float64(rtm.Alloc))
	m.Storage.SetGauge("Frees", float64(rtm.Frees))
	m.Storage.SetGauge("GCCPUFraction", rtm.GCCPUFraction)
	m.Storage.SetGauge("GCSys", float64(rtm.GCSys))
	m.Storage.SetGauge("HeapAlloc", float64(rtm.HeapAlloc))
	m.Storage.SetGauge("HeapIdle", float64(rtm.HeapIdle))
	m.Storage.SetGauge("HeapInuse", float64(rtm.HeapInuse))
	m.Storage.SetGauge("HeapObjects", float64(rtm.HeapObjects))
	m.Storage.SetGauge("HeapReleased", float64(rtm.HeapReleased))
	m.Storage.SetGauge("HeapSys", float64(rtm.HeapSys))
	m.Storage.SetGauge("LastGC", float64(rtm.LastGC))
	m.Storage.SetGauge("Lookups", float64(rtm.Lookups))
	m.Storage.SetGauge("MCacheInuse", float64(rtm.MCacheInuse))
	m.Storage.SetGauge("MCacheSys", float64(rtm.MCacheSys))
	m.Storage.SetGauge("MSpanInuse", float64(rtm.MSpanInuse))
	m.Storage.SetGauge("MSpanSys", float64(rtm.MSpanSys))
	m.Storage.SetGauge("Mallocs", float64(rtm.Mallocs))
	m.Storage.SetGauge("NextGC", float64(rtm.NextGC))
	m.Storage.SetGauge("NumForcedGC", float64(rtm.NumForcedGC))
	m.Storage.SetGauge("NumGC", float64(rtm.NumGC))
	m.Storage.SetGauge("OtherSys", float64(rtm.OtherSys))
	m.Storage.SetGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	m.Storage.SetGauge("StackInuse", float64(rtm.StackInuse))
	m.Storage.SetGauge("StackSys", float64(rtm.StackSys))
	m.Storage.SetGauge("Sys", float64(rtm.Sys))
	m.Storage.SetGauge("TotalAlloc", float64(rtm.TotalAlloc))
}

func (m *runtimeMonitor) collectRuntimeSystem() {
	m.Storage.SetGauge("RandomValue", 1000*rand.Float64())
	m.Storage.AddCounter("PollCount", 1)
}
