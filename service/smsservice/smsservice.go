package main

import (
	_ "embed"
	"zservice/internal/etcdservice"
	"zservice/internal/grpcservice"
	"zservice/service/smsservice/internal"
	"zservice/service/smsservice/smsservice_pb"
	"zservice/zservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

//go:embed version
var Version string

func init() {

	zservice.Init(&zservice.ZServiceConfig{
		Name:    "smsservice",
		Version: Version,
	})

}

func main() {

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			// do something
		},
	})

	grpcS := grpcservice.NewGrpcService(&grpcservice.GrpcServiceConfig{
		Addr:       zservice.Getenv("GRPC_ADDR"),
		EtcdServer: etcdS.Etcd,
		OnStart: func(grpc *grpc.Server) {
			smsservice_pb.RegisterSmsserviceServer(grpc, internal.NewSmsserviceServer())
		},
	})

	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grpcS.ZService)

	grpcS.AddDependService(etcdS.ZService)

	zservice.Start()
	zservice.WaitStart()
	zservice.WaitStop()
}
