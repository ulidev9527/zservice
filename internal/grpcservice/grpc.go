package grpcservice

import (
	"fmt"
	"net"
	"time"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type GrpcService struct {
	*zservice.ZService
	Grpc *grpc.Server
}

type GrpcServiceConfig struct {
	Name    string // 服务名
	Addr    string // 监听地址
	Etcd    *clientv3.Client
	OnStart func(*grpc.Server) // 启动的回调
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

		lis, e := net.Listen("tcp", c.Addr)
		if e != nil {
			s.LogPanic(e)
		}

		// 创建 etcd 客户端
		mgrTarget := fmt.Sprintf("%s/services/%s", zservice.S_GrpcRootKey, zservice.GetName())
		mgr, e := endpoints.NewManager(c.Etcd, mgrTarget)
		if e != nil {
			s.LogPanic(e)
		}

		// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
		lease, e := c.Etcd.Grant(c.Etcd.Ctx(), 10)
		if e != nil {
			s.LogPanic(e)
		}

		list, e := mgr.List(c.Etcd.Ctx())
		if e != nil {
			s.LogPanic(e)
		}
		endpointKey := fmt.Sprintf("%s/%d", mgrTarget, len(list))
		s.LogInfo("grcp endpointKey:", endpointKey)
		// 添加注册节点到 etcd 中，并且携带上租约 id
		// 以 serverName/serverAddr 为 key，serverAddr 为 value
		// serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可
		e = mgr.AddEndpoint(c.Etcd.Ctx(), endpointKey, endpoints.Endpoint{Addr: c.Addr}, clientv3.WithLease(lease.ID))
		if e != nil {
			s.LogPanic(e)
		}

		gs.Grpc = grpc.NewServer(
			grpc.ChainUnaryInterceptor(UnaryServerInterceptor),
			grpc.ChainStreamInterceptor(StreamServerInterceptor),
		)

		if c.OnStart != nil {
			c.OnStart(gs.Grpc)
		}

		// 每隔 5 s进行一次延续租约的动作
		go func() {
			for {
				select {
				case <-time.After(5 * time.Second):
					retryCount := 0
					// 续约操作
					_, e := c.Etcd.KeepAliveOnce(c.Etcd.Ctx(), lease.ID)
					if e != nil {
						retryCount++
						if retryCount > 3 {
							s.LogPanic(e)
						} else {
							s.LogError(e)
						}
					} else {
						retryCount = 0
					}
				case <-c.Etcd.Ctx().Done():
					s.LogPanic(c.Etcd.Ctx().Err())
					return
				}
			}
		}()

		// 启动服务
		go func() {
			s.LogInfof("grpcService listen on %v", c.Addr)
			e := gs.Grpc.Serve(lis)
			if e != nil {
				s.LogPanic(e)
			}
		}()

		go func() {
			s.StartDone()
		}()
	})

	return gs

}
