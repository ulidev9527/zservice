package grpcservice

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ulidev9527/zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type GrpcService struct {
	*zservice.ZService
	GrpcServer *grpc.Server
}

type GrpcServiceConfig struct {
	ListenPort string // 监听端口
	EtcdClient *clientv3.Client
	OnStart    func(*GrpcService) // 启动的回调
}

func NewGrpcService(c *GrpcServiceConfig) *GrpcService {

	if c == nil {
		zservice.LogPanic("GrpcServiceConfig is nil")
		return nil
	}

	name := fmt.Sprint("GrpcService-", c.ListenPort)

	gs := &GrpcService{}
	gs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		// https://ayang.ink/分布式_grpc-基于-etcd-的服务发现/#grpc-服务端

		gs.GrpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(ServerUnaryInterceptor),
			grpc.ChainStreamInterceptor(ServerStreamInterceptor),
		)

		chanConn := make(chan any)
		go func() {

			if c.EtcdClient == nil {
				close(chanConn)
				return
			}

			// 创建 etcd 客户端
			mgrTarget := fmt.Sprintf(S_ServiceName, zservice.GetServiceName())
			mgr, e := endpoints.NewManager(c.EtcdClient, mgrTarget)
			if e != nil {
				s.LogPanic(e)
			}

			isCloseChanconn := false

			for {
				// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
				lease, e := c.EtcdClient.Grant(c.EtcdClient.Ctx(), 10)
				if e != nil {
					s.LogError(e)
					time.Sleep(time.Second) // 等待1秒重连
					continue
				}

				hostName, e := os.Hostname()
				if e != nil {
					s.LogError(e)
					time.Sleep(time.Second) // 等待1秒重连
					continue
				}

				listener := fmt.Sprint(hostName, ":", c.ListenPort)
				endpointKey := fmt.Sprintf("%s/%s", mgrTarget, listener)
				s.LogInfo("grcp endpointKey:", endpointKey)
				// 添加注册节点到 etcd 中，并且携带上租约 id
				// 以 serverName/serverAddr 为 key，serverAddr 为 value
				// serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可

				e = mgr.AddEndpoint(c.EtcdClient.Ctx(), endpointKey, endpoints.Endpoint{Addr: listener}, clientv3.WithLease(lease.ID))
				if e != nil {
					s.LogError(e)
					time.Sleep(time.Second) // 等待1秒重连
					continue
				}

				if !isCloseChanconn { // 控制顺序
					isCloseChanconn = true
					close(chanConn)
				}

				// 处理租约续期，如果续租失败或者租约过期则退出
				for {
					isTimeout := false
					select {
					case <-time.After(5 * time.Second):
						// 租约
						_, err := c.EtcdClient.KeepAliveOnce(context.Background(), lease.ID)
						if err != nil {
							s.LogErrorf("Failed to keep lease alive: %s\n", err.Error())
							isTimeout = true
						}
					case <-c.EtcdClient.Ctx().Done():
						s.LogError(c.EtcdClient.Ctx().Err())
						isTimeout = true
					}
					if isTimeout { // 重连
						break
					}
				}
				s.LogWarn("GRPC Reconnecting...")
			}
		}()

		// 启动 grpc 服务
		s.LogInfof("grpcService listen on :%v", c.ListenPort)
		go func() {
			<-chanConn

			lis, e := net.Listen("tcp", fmt.Sprint(":", c.ListenPort))
			if e != nil {
				s.LogPanic(e)
			}

			if e := gs.GrpcServer.Serve(lis); e != nil {
				s.LogPanic(e)
			}
		}()

		go func() {
			<-chanConn
			s.StartDone()
		}()

		if c.OnStart != nil {
			c.OnStart(gs)
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

	zctx := func() *zservice.Context {
		S2SArr := md.Get(zservice.S_S2S_CTX)
		if len(S2SArr) > 0 {
			zservice.LogDebug(zservice.S_S2S_CTX, S2SArr[0])
			return zservice.NewContext(S2SArr[0])
		} else {
			return zservice.NewContext()
		}

	}()
	zctx.ContextS2S.RequestIP = ipaddr

	// 异常捕获
	defer func() {
		e := recover()
		if e != nil {
			buf := make([]byte, 1<<10)
			stackSize := runtime.Stack(buf, true)
			zctx.LogErrorf("GRPC %v %v :Q %v :E %v :ST %v", ipaddr, info.FullMethod, req, e, string(buf[:stackSize]))
		}
	}()

	resp, e := handler(zctx, req)

	// 打印日志
	if e != nil {
		zctx.LogErrorf("GRPC %v %v :Q %v :E %v", ipaddr, info.FullMethod, req, e)
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
