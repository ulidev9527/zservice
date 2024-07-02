package grpcservice

import (
	"context"
	"fmt"
	"runtime"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcClientOption struct {
	ServiceName string           // 服务名
	EtcdClient  *clientv3.Client // etcd 客户端
	Addr        string           // grcp 服务器地址
}

// GRPC 客户端创建，优先使用 Addr 配置的地址直连，如果 Addr 未配置，使用 ETCD 服务发现
func NewGrpcClient(opt *GrpcClientOption) (*grpc.ClientConn, error) {

	// 配置检查

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(ClientUnaryInterceptor),
	}

	// 直连
	if opt.Addr != "" {
		return grpc.Dial(opt.Addr, grpcOpts...)
	}

	// etcd
	if opt.ServiceName == "" {
		return nil, zservice.NewError("ServiceName is nil")
	}
	if opt.EtcdClient == nil {
		return nil, zservice.NewError("EtcdClient is nil")
	}

	serviceName := fmt.Sprintf(S_ServiceName, opt.ServiceName)
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	builder, e := resolver.NewBuilder(opt.EtcdClient)
	if e != nil {
		return nil, e
	}

	// etcd 需要的内容
	grpcOpts = append(grpcOpts,
		// 注入 etcd resolver
		grpc.WithResolvers(builder),
		// 声明使用的负载均衡策略为 roundrobin，轮询。（测试 target 时去除该注释）
		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)

	// 创建 grpc 连接代理
	conn, e := grpc.Dial(fmt.Sprintf("etcd:///%s", serviceName), grpcOpts...)
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
	zctx := ctx.(*zservice.Context)
	if zctx == nil {
		zctx = zservice.NewContext()
	}

	// 配置metadata
	ctx = metadata.AppendToOutgoingContext(ctx, zservice.S_S2S, zservice.JsonMustMarshalString(zctx.ContextS2S))
	zctx.LogDebug(method, zservice.S_C2S, zservice.JsonMustMarshalString(zctx.ContextS2S))

	// panic
	defer func() {
		e := recover()
		if e != nil {
			buf := make([]byte, 1<<10)
			stackSize := runtime.Stack(buf, true)
			zctx.LogErrorf("GRPC %s :Q %v :E %v %v", method, req, e, string(buf[:stackSize]))
		}
	}()

	// pre-processing
	e := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing

	if e != nil {
		zctx.LogErrorf("GRPC %s :Q %v :E %v", method, req, e)
	} else {
		zctx.LogInfof("GRPC %s :Q %v :S %v", method, req, reply)
	}

	return e
}
