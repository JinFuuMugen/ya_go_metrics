package monitors

import (
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

type nopSender struct{}

func (n nopSender) Process([]storage.Counter, []storage.Gauge) error { return nil }
func (n nopSender) Compress(b []byte) ([]byte, error)                { return b, nil }

func BenchmarkRuntimeCollect(b *testing.B) {
	s := storage.NewStorage()
	m := NewRuntimeMonitor(s, nopSender{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Collect()
	}
}
