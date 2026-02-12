package network

import (
	"context"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSubnetUnaryInterceptor_EmptySubnet_Allows(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	called := false
	handler := func(ctx context.Context, req any) (any, error) {
		called = true
		return "ok", nil
	}

	resp, err := ic(context.Background(), "req", info, handler)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.True(t, called)
}

func TestSubnetUnaryInterceptor_InvalidSubnet_ReturnsInternal(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("not-a-cidr")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := ic(context.Background(), "req", info, handler)
	require.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, st.Code())
}

func TestSubnetUnaryInterceptor_MissingMetadata_PermissionDenied(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("10.0.0.0/8")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := ic(context.Background(), "req", info, handler)
	require.Nil(t, resp)
	require.Error(t, err)

	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestSubnetUnaryInterceptor_MissingXRealIP_PermissionDenied(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("10.0.0.0/8")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("something", "else"))

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := ic(ctx, "req", info, handler)
	require.Nil(t, resp)
	require.Error(t, err)

	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestSubnetUnaryInterceptor_InvalidIP_PermissionDenied(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("10.0.0.0/8")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(HeaderXRealIP, "nope"))

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := ic(ctx, "req", info, handler)
	require.Nil(t, resp)
	require.Error(t, err)

	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestSubnetUnaryInterceptor_IPNotAllowed_PermissionDenied(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("10.0.0.0/8")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(HeaderXRealIP, "192.168.1.10"))

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := ic(ctx, "req", info, handler)
	require.Nil(t, resp)
	require.Error(t, err)

	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestSubnetUnaryInterceptor_IPAllowed_Allows(t *testing.T) {
	_ = logger.Init()

	ic := SubnetUnaryInterceptor("10.0.0.0/8")
	info := &grpc.UnaryServerInfo{FullMethod: "/metrics.Metrics/UpdateMetrics"}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(HeaderXRealIP, "10.1.2.3"))

	called := false
	handler := func(ctx context.Context, req any) (any, error) {
		called = true
		return "ok", nil
	}

	resp, err := ic(ctx, "req", info, handler)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.True(t, called)
}
