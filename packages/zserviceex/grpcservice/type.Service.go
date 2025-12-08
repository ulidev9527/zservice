package grpcservice

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"zserviceapps/packages/zservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// 服务端配置
// type ServiceOption struct {
// 	OnGetEtcd    func() *etcdservice.Service // 如何要使用 etcd 服务发现,这个不能为空
// 	EtcdKey      string                      // * etcd key, 如果 / 开头直接使用, 非 / 开头会添加 /zserviceapps/XXXX
// 	Port         string                      // * 监听端口
// 	GrpcRegister func(*grpc.Server)          // grpc 服务注册
// 	OnStart      func(*Service)              // 启动的回调
// }

type WithXXX func(ser *Service)

func WithPort(port int) WithXXX    { return func(ser *Service) { ser.port = port } }
func WithName(name string) WithXXX { return func(ser *Service) { ser.name = name } }

type OnStart func(ser *Service)

func WithOnStart(onStart OnStart) WithXXX { return func(ser *Service) { ser.onStart = onStart } }

func WithGRPCServer(fn func(grpcSer *grpc.Server)) WithXXX {
	return func(ser *Service) { ser.onRegistRPCServer = fn }
}

type Service struct {
	zservice   *zservice.ZService
	grpcServer *grpc.Server

	port              int                       // 端口
	name              string                    // 服务名称
	onStart           OnStart                   // 启动
	onRegistRPCServer func(rpcSer *grpc.Server) // 注册 RPC 服务
}

var open_log_info = false

func init() {

	zservice.WatchEnvChange("grpcservice_open_log_info", func(key, newVal, oldVal string) {
		open_log_info = zservice.GetenvBool(newVal)
	})

}

func NewService(opts ...WithXXX) *Service {

	ser := &Service{}
	for _, opt := range opts {
		opt(ser)
	}

	if ser.name == "" {
		ser.name = zservice.RandomMD5_XID()
	}
	ser.zservice = zservice.NewService(zservice.ServiceOptions{
		Name: fmt.Sprintf("[grpcservice-%s]", ser.name),
		OnStart: func(_ *zservice.ZService) {

			// https://ayang.ink/分布式_grpc-基于-etcd-的服务发现/#grpc-服务端

			ser.grpcServer = grpc.NewServer(
				grpc.ChainUnaryInterceptor(ServerUnaryInterceptor),
				grpc.ChainStreamInterceptor(ServerStreamInterceptor),
			)

			// 注册 rpc 服务，必须在服务启动前注册
			if ser.onRegistRPCServer != nil {
				ser.onRegistRPCServer(ser.grpcServer)
			}

			// 启动 grpc 服务
			zservice.Go(func() {
				if ser.port == 0 {
					ser.port = zservice.GetFreePort()
				}
				lis, e := net.Listen("tcp", fmt.Sprint(":", ser.port))
				if e != nil {
					ser.zservice.LogError(e)
					os.Exit(1)
				}

				if e := ser.grpcServer.Serve(lis); e != nil {
					ser.zservice.LogError(e)
					os.Exit(1)
				}
			})

			// 等待启动完成
			loopCount := 0
			for {
				if loopCount >= 10 {
					zservice.LogErrorf("can`t start grpc service %s, %s:%d", ser.name, zservice.GetHostname(), ser.port)
					os.Exit(1)
				}

				if zservice.IsPortOpen(zservice.GetHostname(), ser.port, time.Second*3) {
					break
				}
				loopCount++
			}

			if ser.onStart != nil {
				ser.onStart(ser)
			}

		},
	})

	return ser
}

func (ser *Service) GetPort() int { return ser.port }

func ServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// 获取客户端ID
	pr, _ := peer.FromContext(ctx)
	ipaddr := pr.Addr.String()

	// 获取 zservice.Context 和 Trace数据
	md, _ := metadata.FromIncomingContext(ctx)

	var zctx *zservice.Context

	if arr := md.Get("ctx"); len(arr) > 0 {
		zctx = zservice.NewContext(arr[0])
	} else {
		zctx = zservice.NewContext()
	}

	// 异常捕获
	defer func() {
		e := recover()
		if e != nil {
			buf := make([]byte, 1<<10)
			stackSize := runtime.Stack(buf, false)
			zctx.LogErrorf("GRPC %v %v :E %v :ST %v", ipaddr, info.FullMethod, e, string(buf[:stackSize]))
		}
	}()

	resp, e := handler(zctx, req)

	// 打印日志
	if e != nil {
		zctx.LogErrorf("GRPC %v %v :E %v", ipaddr, info.FullMethod, e)
	} else if open_log_info {
		zctx.LogInfof("GRPC %v %v", ipaddr, info.FullMethod)
	}

	return resp, e
}

// ServerStreamInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func ServerStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	return err
}
func (ser *Service) GetZService() *zservice.ZService { return ser.zservice }
