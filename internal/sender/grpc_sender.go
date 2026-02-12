package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/network"
	pb "github.com/JinFuuMugen/ya_go_metrics/internal/proto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type grpcSender struct {
	addr   string
	client pb.MetricsClient
	conn   *grpc.ClientConn
}

// NewGRPCSender creates a new GRPCSender instance using the provided configuration.
func NewGRPCSender(cfg config.AgentConfig) (*grpcSender, error) {
	conn, err := grpc.NewClient(
		cfg.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	return &grpcSender{
		addr:   cfg.GRPCAddr,
		conn:   conn,
		client: pb.NewMetricsClient(conn),
	}, nil
}

func (s *grpcSender) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *grpcSender) Process(counters []storage.Counter, gauges []storage.Gauge) error {
	req := &pb.UpdateMetricsRequest{
		Metrics: make([]*pb.Metric, 0, len(counters)+len(gauges)),
	}

	for _, c := range counters {
		delta := c.GetValue().(int64)
		req.Metrics = append(req.Metrics, &pb.Metric{
			Id:    c.GetName(),
			Type:  pb.Metric_COUNTER,
			Delta: delta,
			Value: 0,
		})
	}

	for _, g := range gauges {
		val := g.GetValue().(float64)
		req.Metrics = append(req.Metrics, &pb.Metric{
			Id:    g.GetName(),
			Type:  pb.Metric_GAUGE,
			Delta: 0,
			Value: val,
		})
	}

	ip, err := network.OutboundIPTo(s.addr)
	if err != nil {
		return fmt.Errorf("cannot determine outbound ip: %w", err)
	}

	md := metadata.New(map[string]string{
		"x-real-ip": ip.String(),
	})

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 5*time.Second)
	defer cancel()

	_, err = s.client.UpdateMetrics(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc UpdateMetrics: %w", err)
	}
	return nil
}

// Compress to fullfill Sender interface
func (s *grpcSender) Compress(data []byte) ([]byte, error) {
	return data, nil
}
