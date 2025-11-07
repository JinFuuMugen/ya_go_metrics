package monitors

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

type monitor struct {
	Storage   storage.Storage
	Processor sender.Sender
}

func NewMonitor(s storage.Storage, p sender.Sender) *monitor {
	return &monitor{s, p}
}

func (m *monitor) collectRuntime() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Storage.SetGauge("BuckHashSys", float64(rtm.BuckHashSys))
	m.Storage.SetGauge("Alloc", float64(rtm.Alloc))               //uint64
	m.Storage.SetGauge("Frees", float64(rtm.Frees))               //uint64
	m.Storage.SetGauge("GCCPUFraction", rtm.GCCPUFraction)        //float64
	m.Storage.SetGauge("GCSys", float64(rtm.GCSys))               //uint64
	m.Storage.SetGauge("HeapAlloc", float64(rtm.HeapAlloc))       //uint64
	m.Storage.SetGauge("HeapIdle", float64(rtm.HeapIdle))         //uint64
	m.Storage.SetGauge("HeapInuse", float64(rtm.HeapInuse))       //uint64
	m.Storage.SetGauge("HeapObjects", float64(rtm.HeapObjects))   //uint64
	m.Storage.SetGauge("HeapReleased", float64(rtm.HeapReleased)) //uint64
	m.Storage.SetGauge("HeapSys", float64(rtm.HeapSys))           //uint64
	m.Storage.SetGauge("LastGC", float64(rtm.LastGC))             //uint64
	m.Storage.SetGauge("Lookups", float64(rtm.Lookups))           //uint64
	m.Storage.SetGauge("MCacheInuse", float64(rtm.MCacheInuse))   //uint64
	m.Storage.SetGauge("MCacheSys", float64(rtm.MCacheSys))       //uint64
	m.Storage.SetGauge("MSpanInuse", float64(rtm.MSpanInuse))     //uint64
	m.Storage.SetGauge("MSpanSys", float64(rtm.MSpanSys))         //uint64
	m.Storage.SetGauge("Mallocs", float64(rtm.Mallocs))           //uint64
	m.Storage.SetGauge("NextGC", float64(rtm.NextGC))             //uint64
	m.Storage.SetGauge("NumForcedGC", float64(rtm.NumForcedGC))   //uint32
	m.Storage.SetGauge("NumGC", float64(rtm.NumGC))               //uint32
	m.Storage.SetGauge("OtherSys", float64(rtm.OtherSys))         //uint64
	m.Storage.SetGauge("PauseTotalNs", float64(rtm.PauseTotalNs)) //uint64
	m.Storage.SetGauge("StackInuse", float64(rtm.StackInuse))     //uint64
	m.Storage.SetGauge("StackSys", float64(rtm.StackSys))         //uint64
	m.Storage.SetGauge("Sys", float64(rtm.Sys))                   //uint64
	m.Storage.SetGauge("TotalAlloc", float64(rtm.TotalAlloc))     //uint64s
}

func (m *monitor) collectSystem() {
	m.Storage.SetGauge("RandomValue", 1000*rand.Float64())
	m.Storage.AddCounter("PollCount", 1)
}

func (m *monitor) CollectMetrics() {
	m.collectRuntime()
	m.collectSystem()
}

func (m *monitor) Dump() error {
	for _, c := range m.Storage.GetCounters() {
		err := m.Processor.Process(c)
		if err != nil {
			return fmt.Errorf("error dumping metric: %w", err)
		}
	}
	for _, g := range m.Storage.GetGauges() {
		err := m.Processor.Process(g)
		if err != nil {
			return fmt.Errorf("error dumping metric: %w", err)
		}
	}
	return nil
}
