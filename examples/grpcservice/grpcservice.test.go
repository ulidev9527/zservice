package main

import (
	"github.com/ulidev9527/zservice/zservice"
	"github.com/ulidev9527/zservice/zserviceex/etcdservice"
	"github.com/ulidev9527/zservice/zserviceex/grpcservice"
)

func main() {
	zservice.Init(zservice.ZserviceOption{
		Name: "test",
	})

	etcdS := etcdservice.NewEtcdService(etcdservice.EtcdServiceOption{
		Addr: "127.0.0.1:2379",
	})

	grpcS := grpcservice.NewGrpcService(grpcservice.GrpcServiceOption{
		GrpcPort: "8123",
		OnStart: func(gs *grpcservice.GrpcService) {

		},
	})

	grpcservice.NewGrpcClientConn(grpcservice.GrpcClientConnOption{
		Addr: "127.0.0.1:8213",
	})

	grpcS.AddDependService(etcdS.ZService)

	zservice.Start()

	zservice.WaitStop()

}
