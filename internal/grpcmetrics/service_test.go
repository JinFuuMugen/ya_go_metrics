package grpcmetrics

import (
	"context"
	"testing"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	pb "github.com/JinFuuMugen/ya_go_metrics/internal/proto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type auditObserver struct {
	events []models.AuditEvent
}

func (o *auditObserver) Notify(e models.AuditEvent) error {
	o.events = append(o.events, e)
	return nil
}

func TestService_UpdateMetrics_EmptyRequest(t *testing.T) {
	_ = logger.Init()

	st := storage.NewStorage()
	p := audit.NewPublisher()
	obs := &auditObserver{}
	p.Subscribe(obs)

	svc := New(st, p)

	resp, err := svc.UpdateMetrics(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, obs.events, 0, "no audit events expected for nil request")

	_, err = st.GetGauge("any")
	require.Error(t, err)
	_, err = st.GetCounter("any")
	require.Error(t, err)
}

func TestService_UpdateMetrics_StoresMetricsAndPublishesAuditWithIP(t *testing.T) {
	_ = logger.Init()

	st := storage.NewStorage()
	p := audit.NewPublisher()
	obs := &auditObserver{}
	p.Subscribe(obs)

	svc := New(st, p)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-real-ip", "10.1.2.3"))

	req := &pb.UpdateMetricsRequest{
		Metrics: []*pb.Metric{
			{Id: "Alloc", Type: pb.Metric_GAUGE, Value: 123.5},
			{Id: "PollCount", Type: pb.Metric_COUNTER, Delta: 7},
			nil,
			{Id: "", Type: pb.Metric_GAUGE, Value: 1},
		},
	}

	before := time.Now().Unix()

	resp, err := svc.UpdateMetrics(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	g, err := st.GetGauge("Alloc")
	require.NoError(t, err)
	require.NotNil(t, g)
	require.Equal(t, 123.5, g.GetValue().(float64))

	c, err := st.GetCounter("PollCount")
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, int64(7), c.GetValue().(int64))

	require.Len(t, obs.events, 1)
	ev := obs.events[0]

	require.Equal(t, "10.1.2.3", ev.IPAddress)
	require.ElementsMatch(t, []string{"Alloc", "PollCount"}, ev.Metrics)

	require.GreaterOrEqual(t, ev.TS, before)
	require.LessOrEqual(t, ev.TS, time.Now().Unix())
}

func TestService_UpdateMetrics_NoAuditPublisher(t *testing.T) {
	_ = logger.Init()

	st := storage.NewStorage()
	svc := New(st, nil)

	req := &pb.UpdateMetricsRequest{
		Metrics: []*pb.Metric{
			{Id: "G1", Type: pb.Metric_GAUGE, Value: 1.25},
		},
	}

	resp, err := svc.UpdateMetrics(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	g, err := st.GetGauge("G1")
	require.NoError(t, err)
	require.NotNil(t, g)
	require.Equal(t, 1.25, g.GetValue().(float64))
}

func TestService_UpdateMetrics_IgnoresNilAndEmptyID(t *testing.T) {
	_ = logger.Init()

	st := storage.NewStorage()
	p := audit.NewPublisher()
	obs := &auditObserver{}
	p.Subscribe(obs)

	svc := New(st, p)

	req := &pb.UpdateMetricsRequest{
		Metrics: []*pb.Metric{
			nil,
			{Id: "", Type: pb.Metric_COUNTER, Delta: 10},
		},
	}

	resp, err := svc.UpdateMetrics(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	_, err = st.GetGauge("anything")
	require.Error(t, err)
	_, err = st.GetCounter("anything")
	require.Error(t, err)

	require.Len(t, obs.events, 0)
}
