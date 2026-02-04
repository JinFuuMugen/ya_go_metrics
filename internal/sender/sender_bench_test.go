package sender

import (
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func BenchmarkSenderProcess(b *testing.B) {
	cfg := config.Config{
		Addr: "localhost:8080",
		Key:  "",
	}

	s := NewSender(cfg, nil)

	st := storage.NewStorage()
	for i := 0; i < 100; i++ {
		st.SetGauge("gauge", float64(i))
		st.AddCounter("counter", 1)
	}

	counters := st.GetCounters()
	gauges := st.GetGauges()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Process(counters, gauges)
	}
}
