package storage

import "testing"

func BenchmarkSetGauge(b *testing.B) {
	s := NewStorage()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.SetGauge("Alloc", float64(i))
	}
}

func BenchmarkAddCounter(b *testing.B) {
	s := NewStorage()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.AddCounter("PollCount", 1)
	}
}
