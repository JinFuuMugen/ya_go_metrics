package storage

import (
	"fmt"
	"strconv"
	"strings"
)

// MetricTypeGauge represents the gauge metric type.
const MetricTypeGauge = "gauge"

// MetricTypeCounter represents the counter metric type.
const MetricTypeCounter = "counter"

type (

	// Metric defines a common interface for all metric types.
	Metric interface {

		// GetType returns the metric type.
		GetType() string

		// GetName returns the metric name.
		GetName() string

		// GetValueString returns the metric value as a string.
		GetValueString() string

		// GetValue returns the raw metric value.
		GetValue() interface{}
	}

	// Storage defines an interface for metric storage backends.
	Storage interface {
		SetGauge(string, float64)
		AddCounter(string, int64)
		GetCounters() []Counter
		GetGauges() []Gauge
		GetCounter(string) (Counter, error)
		GetGauge(string) (Gauge, error)
	}

	// Counter represents a counter metric.
	Counter struct {
		Name  string
		Type  string
		Value int64
	}
	Gauge struct {
		Name  string
		Type  string
		Value float64
	}
)

func (c Counter) GetType() string {
	return c.Type
}

func (c Counter) GetName() string {
	return c.Name
}

func (c Counter) GetValue() interface{} {
	return c.Value
}

func (c Counter) GetValueString() string {
	return strconv.FormatInt(c.Value, 10)
}

func (g Gauge) GetType() string {
	return g.Type
}

func (g Gauge) GetName() string {
	return g.Name
}

func (g Gauge) GetValue() interface{} {
	return g.Value
}

func (g Gauge) GetValueString() string {
	f := func(num float64) string {
		s := fmt.Sprintf(`%.4f`, num)
		return strings.TrimRight(strings.TrimRight(s, `0`), `.`)
	}
	return f(g.Value)
}

// NewStorage creates a new in-memory metric storage.
func NewStorage() Storage {
	return &MemStorage{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
	}
}

var defaultStorage = NewStorage()

func GetCounter(k string) (Counter, error) {
	return defaultStorage.GetCounter(k)
}

func GetGauge(k string) (Gauge, error) {
	return defaultStorage.GetGauge(k)
}

func AddCounter(k string, v int64) {
	defaultStorage.AddCounter(k, v)
}

func SetGauge(k string, v float64) {
	defaultStorage.SetGauge(k, v)
}
func GetCounters() []Counter {
	return defaultStorage.GetCounters()
}

func GetGauges() []Gauge {
	return defaultStorage.GetGauges()
}
