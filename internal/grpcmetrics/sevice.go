package grpcmetrics

import (
	"context"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	pb "github.com/JinFuuMugen/ya_go_metrics/internal/proto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	pb.UnimplementedMetricsServer

	st  storage.Storage
	aud *audit.Publisher
}

func New(st storage.Storage, aud *audit.Publisher) *Service {
	return &Service{st: st, aud: aud}
}

func (s *Service) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	if req == nil || len(req.Metrics) == 0 {
		return &pb.UpdateMetricsResponse{}, nil
	}

	affected := make([]string, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		if m == nil || m.Id == "" {
			continue
		}

		switch m.Type {
		case pb.Metric_COUNTER:
			s.st.AddCounter(m.Id, m.Delta)
		case pb.Metric_GAUGE:
			s.st.SetGauge(m.Id, m.Value)
		default:
		}

		affected = append(affected, m.Id)
	}

	if s.aud != nil && len(affected) > 0 {
		ip := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-real-ip"); len(vals) > 0 {
				ip = vals[0]
			}
		}

		s.aud.Publish(models.AuditEvent{
			TS:        time.Now().Unix(),
			Metrics:   affected,
			IPAddress: ip,
		})
	}

	return &pb.UpdateMetricsResponse{}, nil
}
