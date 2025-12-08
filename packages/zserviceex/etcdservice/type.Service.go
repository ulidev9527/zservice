package etcdservice

import (
	"fmt"
	"time"

	"zserviceapps/packages/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Service struct {
	*zservice.ZService
	Client *clientv3.Client
}

type ServiceOption struct {
	Name    string         // 名字， 默认取 etcd_XXX
	Addrs   []string       // ETCD 服务地址
	OnStart func(*Service) // 启动的回调
}

func NewService(c ServiceOption) *Service {

	if c.Name == "" {
		c.Name = fmt.Sprint("etcd_", zservice.RandomXID())
	}

	ser := &Service{}
	ser.ZService = zservice.NewService(zservice.ServiceOptions{
		Name: c.Name,
		OnStart: func(_ *zservice.ZService) {
			for {
				etcd, e := clientv3.New(clientv3.Config{
					Endpoints:   c.Addrs,
					DialTimeout: 5 * time.Second,
				})

				if e != nil {
					ser.LogError("has error, waiting 5s again:", e)
					time.Sleep(time.Second * 5)
					continue
				}
				ser.Client = etcd

				if fails := ser.Ping(ser.GetLauncherCtx()); len(fails) > 0 {
					time.Sleep(time.Second * 5)
					continue
				}

				if c.OnStart != nil {
					c.OnStart(ser)
				}
				break
			}
		},
	})
	return ser
}
