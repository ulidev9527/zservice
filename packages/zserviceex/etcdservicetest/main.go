package main

import (
	"zserviceapps/packages/zservice"
	"zserviceapps/packages/zserviceex/etcdservice"
)

func main() {

	zservice.NewService(zservice.ServiceOptions{
		Name: "etcd_test",
		OnStart: func(z *zservice.ZService) {

		},
	}).AddDependService(
		etcdservice.NewService(etcdservice.ServiceOption{
			Addrs: []string{"127.0.0.1:2379"},
		}).ZService,
	)
	zservice.Start()

	zservice.WaitStart()
}
