package main

import (
	"github.com/ulidev9527/zservice/zservice"
	"github.com/ulidev9527/zservice/zserviceex/etcdservice"
)

func main() {

	zservice.Init(zservice.ZserviceOption{})

	etcdservice.NewEtcdService(etcdservice.EtcdServiceOption{
		Addr: "127.0.0.1:2379",
		OnStart: func(es *etcdservice.EtcdService) {

			es.LogInfo("start")

		},
	})

	zservice.Start()
	zservice.WaitStop()
}
