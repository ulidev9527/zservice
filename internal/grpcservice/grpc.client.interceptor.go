package grpcservice

import (
	"context"
	"encoding/json"
	"zservice/zservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// https://www.lixueduan.com/posts/grpc/05-interceptor/
// https://blog.csdn.net/qq_30614345/article/details/134470773

// unaryInterceptor 一个简单的 unary interceptor 示例。
func ClientUnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 获取 context 中的 zservice.Context
	zctx := GetCtxEX(ctx)
	if zctx == nil {
		zctx = zservice.NewEmptyContext()
	}

	// 配置metadata
	traceJson, _ := json.Marshal(zctx.ContextTrace)
	ctx = metadata.AppendToOutgoingContext(ctx, zservice.S_TraceKey, string(traceJson))

	// panic
	defer func() {
		e := recover()
		if e != nil {
			zctx.LogErrorf("RPC %s :Q %v :E %v", method, req, e)
		}
	}()

	// pre-processing
	e := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing

	if e != nil {
		zctx.LogErrorf("RPC %s :Q %v :E %v", method, req, e)
	} else {
		zctx.LogInfof("RPC %s :Q %v :S %v", method, req, reply)
	}

	return e
}
