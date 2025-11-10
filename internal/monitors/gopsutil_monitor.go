package monitors

import (
	"fmt"
	"log"

	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type gopsutilMonitor struct {
	Storage   storage.Storage
	Processor sender.Sender
}

func NewGopsutilMonitor(s storage.Storage, p sender.Sender) GopsutilMonitor {
	return &gopsutilMonitor{s, p}
}

func (m *gopsutilMonitor) Collect() {
	m.CollectGopsutil()
}

func (m *gopsutilMonitor) CollectGopsutil() {

	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("cannot get memory info: %s", err)
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Fatalf("cannot get CPU info: %s", err)
	}

	m.Storage.SetGauge("TotalMemory", float64(vm.Total))
	m.Storage.SetGauge("FreeMemory", float64(vm.Available))
	m.Storage.SetGauge("CPUutilization1", cpuPercent[0])
}

func (m *gopsutilMonitor) Dump() error {
	c := m.Storage.GetCounters()
	g := m.Storage.GetGauges()
	err := m.Processor.Process(c, g)
	if err != nil {
		return fmt.Errorf("error dumping metric: %w", err)
	}
	return nil
}

func (m *gopsutilMonitor) SetProcessor(p sender.Sender) {
	m.Processor = p
}
