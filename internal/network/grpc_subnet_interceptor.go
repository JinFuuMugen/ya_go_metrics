package network

import (
	"context"
	"net"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const HeaderXRealIP = "x-real-ip"

// SubnetUnaryInterceptor checks if request contains valid x-real-ip header
func SubnetUnaryInterceptor(trustedSubnet string) grpc.UnaryServerInterceptor {
	var (
		ipnet  *net.IPNet
		cfgErr error
	)

	if trustedSubnet != "" {
		_, parsed, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			cfgErr = err
			logger.Errorf("invalid trusted subnet CIDR %q: %v", trustedSubnet, err)
		} else {
			ipnet = parsed
		}
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if cfgErr != nil {
			return nil, status.Error(codes.Internal, "invalid trusted subnet configuration")
		}
		if ipnet == nil {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "missing metadata")
		}

		vals := md.Get(HeaderXRealIP)
		if len(vals) == 0 {
			return nil, status.Error(codes.PermissionDenied, "missing x-real-ip metadata")
		}

		ip := net.ParseIP(vals[0])
		if ip == nil {
			return nil, status.Error(codes.PermissionDenied, "invalid x-real-ip metadata")
		}

		if v4 := ip.To4(); v4 != nil {
			ip = v4
		}

		if !ipnet.Contains(ip) {
			return nil, status.Error(codes.PermissionDenied, "ip not allowed")
		}

		return handler(ctx, req)
	}
}
