package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient
var zauthInitConfig *ZAuthInitConfig

type ZAuthInitConfig struct {
	ServiceName string // 权限服务名称
	Etcd        *clientv3.Client
	GrpcAddr    string // rpc addr
	UseGrpcEtcd bool   // 是否使用 rpc + etcd
}

func Init(c *ZAuthInitConfig) {
	if c == nil {
		zservice.LogPanic("ZAuthInitConfig is nil")
		return
	}
	zauthInitConfig = c
	func() {
		conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
			ServiceName: c.ServiceName,
			EtcdClient:  c.Etcd,
			Addr:        c.GrpcAddr,
			UseEtcd:     c.UseGrpcEtcd,
		})
		if e != nil {
			zservice.LogPanic(e)
			return
		}

		grpcClient = zauth_pb.NewZauthClient(conn)
	}()

	if c.ServiceName == "" {
		zservice.LogPanic("ZauthServiceName is nil")
	}

	// 服务配置改变监听
	go internal.EV_Watch_Config_ServiceFileConfigChange(c.Etcd, zservice.GetServiceName(), func(s string) {
		fileConfigCache.Delete(s)
		zservice.LogInfo("Update config ", s)
	})

}
func GetGrpcClient() zauth_pb.ZauthClient {
	return grpcClient
}
func GetGrpcServiceName() string {
	return zauthInitConfig.ServiceName
}
