package models

import "errors"

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) SetDelta(delta int64) {
	m.Delta = &delta
}
func (m *Metrics) SetValue(value float64) {
	m.Value = &value
}
func (m *Metrics) GetValue() (float64, error) {
	if m.Value == nil {
		return 0, errors.New("no value")
	}
	return *m.Value, nil
}

func (m *Metrics) GetDelta() (int64, error) {
	if m.Delta == nil {
		return 0, errors.New("no delta")
	}
	return *m.Delta, nil
}
