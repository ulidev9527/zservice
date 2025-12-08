package grpcservice

import (
	"context"
	"fmt"
	"runtime"

	"zserviceapps/packages/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// 客户端配置
type GrpcClientConnOption struct {
	// addr 和 etcd 二选一
	EtcdKey    string           // * etcd key
	EtcdClient *clientv3.Client // etcd 客户端
	Addr       string           // grcp 服务器地址
}

// GRPC 客户端创建，优先使用 Addr 配置的地址直连，如果 Addr 未配置，使用 ETCD 服务发现
func NewGrpcClientConn(ctx *zservice.Context, opt GrpcClientConnOption) (*grpc.ClientConn, error) {

	grpcDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(ClientUnaryInterceptor),
		grpc.WithBlock(),
	}

	// 直连
	if opt.Addr != "" {
		ctx.LogInfof("grpc dial: %s", opt.Addr)
		return grpc.Dial(opt.Addr, grpcDialOpts...)
	}

	if opt.EtcdClient == nil {
		return nil, zservice.NewError("EtcdClient is nil")
	}

	// etcd
	if opt.EtcdKey == "" {
		return nil, zservice.NewError("EtcdKey is nil")
	}

	if opt.EtcdKey[0] != '/' {
		opt.EtcdKey = "/zserviceapps/" + opt.EtcdKey
	}

	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	builder, e := resolver.NewBuilder(opt.EtcdClient)
	if e != nil {
		return nil, e
	}

	// etcd 需要的内容
	grpcDialOpts = append(grpcDialOpts,
		// 注入 etcd resolver
		grpc.WithResolvers(builder),
		// 声明使用的负载均衡策略为 roundrobin，轮询。（测试 target 时去除该注释）
		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)

	ctx.LogInfof("grpc dial: %s", opt.EtcdKey)

	// 创建 grpc 连接代理
	conn, e := grpc.NewClient(fmt.Sprintf("etcd:///%s", opt.EtcdKey), grpcDialOpts...)
	if e != nil {
		return nil, e
	}

	return conn, nil
}

// DirectConnectGRPC 直连 grpc 服务，无需 TLS
func DirectConnectGRPC(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if len(opts) == 0 {
		opts = append(opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(ClientUnaryInterceptor),
			grpc.WithBlock())
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
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
	ctx = metadata.AppendToOutgoingContext(ctx, "ctx", zctx.ToContextString())

	// panic
	defer func() {
		e := recover()
		if e != nil {
			buf := make([]byte, 1<<10)
			stackSize := runtime.Stack(buf, true)
			zctx.LogErrorf("GRPC %s :E %v %v", method, e, string(buf[:stackSize]))
		}
	}()

	// pre-processing
	e := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing

	if e != nil {
		zctx.LogErrorf("GRPC %s :E %v", method, e)
	} else if open_log_info {
		zctx.LogInfof("GRPC %s", method)
	}

	return e
}
