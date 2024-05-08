package grpcservice

import (
	"context"
	"log"
	"zservice/zservice"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("before handling. Info: %+v", info)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)

	zservice.LogPanic()

	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("before handling. Info: %+v", info)
	err := handler(srv, ss)
	log.Printf("after handling. err: %v", err)
	return err
}
