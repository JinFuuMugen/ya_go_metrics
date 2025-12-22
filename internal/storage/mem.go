package storage

import (
	"errors"
)

// MemStorage is an in-memory storage for metrics.
type MemStorage struct {
	// GaugeMap stores gauge metrics by name.
	GaugeMap map[string]float64

	// CounterMap stores counter metrics by name.
	CounterMap map[string]int64
}

// SetGauge sets the value of a gauge metric.
func (ms *MemStorage) SetGauge(key string, value float64) {
	ms.GaugeMap[key] = value
}

// AddCounter increments the value of a counter metric.
func (ms *MemStorage) AddCounter(key string, value int64) {
	_, keyExists := ms.CounterMap[key]
	if keyExists {
		ms.CounterMap[key] += value
	} else {
		ms.CounterMap[key] = value
	}
}

// GetGauges returns all stored gauge metrics.
func (ms *MemStorage) GetGauges() []Gauge {
	var gauges []Gauge
	for k, v := range ms.GaugeMap {
		gauges = append(gauges, Gauge{Name: k, Type: MetricTypeGauge, Value: v})
	}
	return gauges
}

// GetCounters returns all stored counter metrics.
func (ms *MemStorage) GetCounters() []Counter {
	var counters []Counter
	for k, v := range ms.CounterMap {
		counters = append(counters, Counter{Name: k, Type: MetricTypeCounter, Value: v})
	}
	return counters
}

// GetCounter returns a counter metric by name.
func (ms *MemStorage) GetCounter(k string) (Counter, error) {
	c, exists := ms.CounterMap[k]
	if exists {
		return Counter{Name: k, Type: MetricTypeCounter, Value: c}, nil
	} else {
		return Counter{}, errors.New("missing key")
	}
}

// GetGauge returns a gauge metric by name.
func (ms *MemStorage) GetGauge(k string) (Gauge, error) {
	g, exists := ms.GaugeMap[k]
	if exists {
		return Gauge{Name: k, Type: MetricTypeGauge, Value: g}, nil
	} else {
		return Gauge{}, errors.New("missing key")
	}
}
