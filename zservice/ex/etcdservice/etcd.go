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
	Name    string                 // 服务名
	Addr    string                 // 服务地址
	OnStart func(*clientv3.Client) // 启动的回调
}

func NewEtcdService(c *EtcdServiceConfig) *EtcdService {

	if c == nil {
		zservice.LogPanic("EtcdServiceConfig is nil")
		return nil
	}

	name := "EtcdService"

	if c.Name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	es := &EtcdService{}
	es.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		es.LogInfof("etcdService listen on %v", c.Addr)

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, e := es.Etcd.Status(timeoutCtx, c.Addr)
		if e != nil {
			s.LogPanic(e)
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