package grpcservice

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type GrpcService struct {
	*zservice.ZService
	Server *grpc.Server
}

type GrpcServiceConfig struct {
	Name       string // 服务名
	ListenAddr string // 监听地址
	EtcdServer *clientv3.Client
	OnStart    func(*grpc.Server) // 启 动的回调
}

func NewGrpcService(c *GrpcServiceConfig) *GrpcService {

	if c == nil {
		zservice.LogPanic("GrpcServiceConfig is nil")
		return nil
	}

	name := "GrpcService"

	if c.Name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	gs := &GrpcService{}
	gs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		// https://ayang.ink/分布式_grpc-基于-etcd-的服务发现/#grpc-服务端

		lis, e := net.Listen("tcp", c.ListenAddr)
		if e != nil {
			s.LogPanic(e)
		}

		gs.Server = grpc.NewServer(
			grpc.ChainUnaryInterceptor(ServerUnaryInterceptor),
			grpc.ChainStreamInterceptor(ServerStreamInterceptor),
		)

		// 创建 etcd 客户端
		mgrTarget := fmt.Sprintf(S_ServiceName, zservice.GetServiceName())
		mgr, e := endpoints.NewManager(c.EtcdServer, mgrTarget)
		if e != nil {
			s.LogPanic(e)
		}

		go func() {
			reConnCount := 0 // 重连次数
			for {
				// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
				lease, e := c.EtcdServer.Grant(c.EtcdServer.Ctx(), 10)
				if e != nil {
					s.LogPanic(e)
				}

				if ips, e := zservice.GetIp(); e != nil {
					s.LogPanic(e)
				} else {

					port := ""
					if strings.Contains(c.ListenAddr, ":") {
						port = strings.Split(c.ListenAddr, ":")[1]
						port = ":" + port
					}
					for _, ipaddr := range ips {

						listener := ipaddr + port
						endpointKey := fmt.Sprintf("%s/%s", mgrTarget, listener)
						s.LogInfo("grcp endpointKey:", endpointKey)
						// 添加注册节点到 etcd 中，并且携带上租约 id
						// 以 serverName/serverAddr 为 key，serverAddr 为 value
						// serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可

						e := mgr.AddEndpoint(c.EtcdServer.Ctx(), endpointKey, endpoints.Endpoint{Addr: listener}, clientv3.WithLease(lease.ID))
						if e != nil {
							s.LogPanic(e)
						}
					}
				}

				// 处理租约续期，如果续租失败或者租约过期则退出
				for {
					isTimeout := false
					select {
					case <-time.After(5 * time.Second):
						// 租约
						_, err := c.EtcdServer.KeepAliveOnce(context.Background(), lease.ID)
						if err != nil {
							fmt.Printf("Failed to keep lease alive: %s\n", err.Error())
							isTimeout = true
						}
					case <-c.EtcdServer.Ctx().Done():
						s.LogPanic(c.EtcdServer.Ctx().Err())
					}
					if isTimeout { // 超时重连
						break
					}
				}
				time.Sleep(1 * time.Second) // 等待1秒重连

				reConnCount++
				if reConnCount > 10 {
					s.LogPanic("GRPC connect failed!")
				} else {
					s.LogWarn("GRPC Reconnecting...")
				}
			}
		}()

		// 启动 grpc 服务
		go func() {
			s.LogInfof("grpcService listen on %v", c.ListenAddr)
			e := gs.Server.Serve(lis)
			if e != nil {
				s.LogPanic(e)
			}
		}()

		go func() {
			s.StartDone()
		}()

		if c.OnStart != nil {
			c.OnStart(gs.Server)
		}

	})

	return gs
}

func ServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// 获取客户端ID
	pr, _ := peer.FromContext(ctx)
	ipaddr := strings.Split(pr.Addr.String(), ":")[0]

	// 获取 zservice.Context 和 Trace数据
	md, _ := metadata.FromIncomingContext(ctx)
	zctx := zservice.NewContext(md.Get(zservice.S_S2S)[0])
	zctx.ContextS2S.RequestIP = ipaddr
	ctx = context.WithValue(ctx, GRPC_contextEX_Middleware_Key, zctx)

	// 异常捕获
	defer func() {
		e := recover()
		if e != nil {
			buf := make([]byte, 1<<10)
			stackSize := runtime.Stack(buf, true)
			zctx.LogErrorf("GRPC %v %v :Q %v :E %v :ST %v", ipaddr, info.FullMethod, req, e, string(buf[:stackSize]))
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
