package grpcservice

import (
	"context"
	"strings"
	"zservice/zservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func ServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// 获取 zservice.Context 和 Trace数据
	md, _ := metadata.FromIncomingContext(ctx)
	zctx := zservice.NewContext(md.Get(zservice.S_TraceKey)[0])
	ctx = context.WithValue(ctx, GRPC_contextEX_Middleware_Key, zctx)

	// 获取客户端ID
	pr, _ := peer.FromContext(ctx)
	ipaddr := strings.Split(pr.Addr.String(), ":")[0]

	// 异常捕获
	defer func() {
		e := recover()
		if e != nil {
			zctx.LogError(e, "GRPC %v %v :Q %v :E %v", ipaddr, info.FullMethod, req, e)
		}
	}()

	resp, e := handler(ctx, req)

	// 打印日志
	if e != nil {
		zctx.LogError(e, "GRPC %v %v :Q %v :E %v", ipaddr, info.FullMethod, req, e)
	} else {
		zctx.LogInfof("GRPC %v %v :Q %v :S %v", ipaddr, info.FullMethod, req, resp)
	}

	return resp, e
}

// ServerStreamInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func ServerStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	return err
}
