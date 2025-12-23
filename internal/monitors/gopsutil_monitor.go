package monitors

import (
	"fmt"

	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type gopsutilMonitor struct {
	Storage   storage.Storage
	Processor sender.Sender
}

// NewGopsutilMonitor creates a new gopsutil-based monitor.
func NewGopsutilMonitor(s storage.Storage, p sender.Sender) GopsutilMonitor {
	return &gopsutilMonitor{s, p}
}

// Collect collects system metrics using gopsutil.
func (m *gopsutilMonitor) Collect() error {
	return m.CollectGopsutil()
}

// CollectGopsutil collects CPU and memory metrics from the system.
func (m *gopsutilMonitor) CollectGopsutil() error {

	vm, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("cannot get memory info: %w", err)
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return fmt.Errorf("cannot get CPU info: %w", err)
	}

	m.Storage.SetGauge("TotalMemory", float64(vm.Total))
	m.Storage.SetGauge("FreeMemory", float64(vm.Available))
	m.Storage.SetGauge("CPUutilization1", cpuPercent[0])

	return nil
}

// Dump sends collected metrics using the configured Sender.
func (m *gopsutilMonitor) Dump() error {
	c := m.Storage.GetCounters()
	g := m.Storage.GetGauges()
	err := m.Processor.Process(c, g)
	if err != nil {
		return fmt.Errorf("error dumping metric: %w", err)
	}
	return nil
}

// SetProcessor updates the Sender used for dumping metrics.
func (m *gopsutilMonitor) SetProcessor(p sender.Sender) {
	m.Processor = p
}
