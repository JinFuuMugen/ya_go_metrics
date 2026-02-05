package models

import "errors"

var errNoDelta error = errors.New("no delta")
var errNoValue error = errors.New("no value")

// Metrics represents a metric entity.
//
//generate:reset
//go:generate go run ../../cmd/reset/main.go
type Metrics struct {
	// ID is the metric name.
	ID string `json:"id"`

	// MType specifies the metric type
	MType string `json:"type"`

	// Delta stores the value for counter metrics. It is omitted for gauge metrics.
	Delta *int64 `json:"delta,omitempty"`

	// Value stores the value for gauge metrics. It is omitted for counter metrics.
	Value *float64 `json:"value,omitempty"`
}

// SetDelta sets the delta value for a counter metric.
func (m *Metrics) SetDelta(delta int64) {
	m.Delta = &delta
}

// SetValue sets the value for a gauge metric.
func (m *Metrics) SetValue(value float64) {
	m.Value = &value
}

// GetValue returns the gauge metric value.
// An error is returned if the metric does not contain a value.
func (m *Metrics) GetValue() (float64, error) {
	if m.Value == nil {
		return 0, errNoValue
	}
	return *m.Value, nil
}

// GetDelta returns the counter metric delta.
// An error is returned if the metric does not contain a delta.
func (m *Metrics) GetDelta() (int64, error) {
	if m.Delta == nil {
		return 0, errNoDelta
	}
	return *m.Delta, nil
}
