package etcdservice

import (
	"context"
	"fmt"
	"time"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdService struct {
	*zservice.ZService
	Etcd *clientv3.Client
}

type EtcdServiceConfig struct {
	Addr    string                 // ETCD 服务地址
	OnStart func(*clientv3.Client) // 启动的回调
}

func NewEtcdService(c *EtcdServiceConfig) *EtcdService {

	if c == nil {
		zservice.LogPanic("EtcdServiceConfig is nil")
		return nil
	}

	name := fmt.Sprint("EtcdService-", c.Addr)

	es := &EtcdService{}
	es.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		s.LogInfof("etcdService listen on %v", c.Addr)

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if status, e := es.Etcd.Status(timeoutCtx, c.Addr); e != nil {
			s.LogPanic(e)
		} else {
			s.LogInfo("ETCD Status:", string(zservice.JsonMustMarshal(status)))
		}
		if c.OnStart != nil {
			c.OnStart(es.Etcd)
		}
		s.StartDone()

	})

	etcd, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{c.Addr},
		DialTimeout: 5 * time.Second,
	})

	if e != nil {
		es.LogPanic(e)
	}

	es.Etcd = etcd
	return es
}
