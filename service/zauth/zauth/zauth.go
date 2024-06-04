package zauth

import (
	"sync"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient
var fileConfigMap = &sync.Map{} // 文件配置映射
var zauthInitConfig *ZAuthInitConfig

type ZAuthInitConfig struct {
	ZauthServiceName string // 权限服务名称
	Etcd             *clientv3.Client
	Redis            *redisservice.GoRedisEX
	NsqConsumerAddrs string // nsq consumer addr
	IsNsqdAddr       bool
}

func Init(c *ZAuthInitConfig) {
	zauthInitConfig = c
	func() {
		conn, e := grpcservice.NewGrpcClient(&grpcservice.GrpcClientConfig{
			ZauthServiceName: c.ZauthServiceName,
			EtcdServer:       c.Etcd,
		})
		if e != nil {
			zservice.LogPanic(e)
			return
		}

		grpcClient = zauth_pb.NewZauthClient(conn)
	}()

	if c.ZauthServiceName == "" {
		zservice.LogPanic("ZauthServiceName is nil")
	}

	// 服务配置改变监听
	go internal.EV_Watch_Config_ServiceFileConfigChange(c.Etcd, zservice.GetServiceName(), func(s string) {
		fileConfigMap.Delete(s)
		zservice.LogInfo("Update config ", s)
	})
}
