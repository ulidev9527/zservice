package grpcservice

import (
	"context"
	"encoding/json"
	"fmt"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcClientConfig struct { // etcd 和 addr 二选一
	EtcdServiceName string           // 服务名
	EtcdServer      *clientv3.Client // etcd 客户端
	Addr            string           // grcp 服务器地址
}

func NewGrpcClient(c *GrpcClientConfig) (*grpc.ClientConn, error) {

	// etcd 和 addr 二选一
	// 直连
	if c.EtcdServer == nil {
		return grpc.Dial("0.0.0.0:3002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// etcd
	serviceName := fmt.Sprintf(S_ServiceName, c.EtcdServiceName)
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	builder, e := resolver.NewBuilder(c.EtcdServer)
	if e != nil {
		return nil, e
	}

	// 创建 grpc 连接代理
	conn, e := grpc.Dial(
		// 服务名称
		fmt.Sprintf("etcd:///%s", serviceName),
		// 注入 etcd resolver
		grpc.WithResolvers(builder),
		// 声明使用的负载均衡策略为 roundrobin，轮询。（测试 target 时去除该注释）
		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(ClientUnaryInterceptor),
	)
	if e != nil {
		return nil, e
	}

	return conn, nil
}

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
