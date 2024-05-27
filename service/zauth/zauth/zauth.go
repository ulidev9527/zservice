package zauth

import (
	"fmt"
	"sync"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/ex/redisservice"

	"github.com/nsqio/go-nsq"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var grpcClient zauth_pb.ZauthClient
var fileConfigMap = &sync.Map{}
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

	nsqservice.NewNsqConsumer(&nsqservice.NsqConsumerConfig{
		Addrs:      c.NsqConsumerAddrs,
		IsNsqdAddr: c.IsNsqdAddr,
		Topic:      internal.NSQ_FileConfig_Change,
		Channel:    fmt.Sprintf("%s-%s", zservice.GetServiceName(), zservice.RandomXID()),
		OnMessage: func(m *nsq.Message) error {
			fileName := string(m.Body)
			zservice.LogInfo("Update config ", fileName)
			fileConfigMap.Delete(fileName)
			return nil
		},
	})

}
