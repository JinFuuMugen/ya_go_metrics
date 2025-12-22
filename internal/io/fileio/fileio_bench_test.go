package fileio

import (
	"os"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func BenchmarkSaveMetricsFile(b *testing.B) {
	tmp, _ := os.CreateTemp("", "metrics-*.json")
	defer os.Remove(tmp.Name())

	st := storage.NewStorage()
	for i := 0; i < 100; i++ {
		st.SetGauge("gauge", float64(i))
		st.AddCounter("counter", 1)
	}

	c := st.GetCounters()
	g := st.GetGauges()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SaveMetricsFile(tmp.Name(), c, g)
	}
}
