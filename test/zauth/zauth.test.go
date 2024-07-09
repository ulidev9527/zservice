package main

import (
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/test/zauth/internal"
	"zservice/zservice"
	"zservice/zserviceex/etcdservice"
)

func init() {

	zservice.Init("zauth.test", "0.1.0")
}

func main() {

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(s *etcdservice.EtcdService) {

		},
	})

	grpcClient := zservice.NewService("zauth.grpc", func(z *zservice.ZService) {
		defer z.StartDone()
		ctx := zservice.NewContext()
		zauth.Init(&zauth.ZAuthInitOption{
			ServiceName: zservice.Getenv("ZAUTH_SERVICENAME"),
			EtcdService: etcdS,
			GrpcAddr:    zservice.Getenv("ZAUTH_GRPCADDR"),
		})

		zauth.ServiceRegist(ctx, &zauth_pb.ServiceRegist_REQ{})
	})
	grpcClient.AddDependService(etcdS.ZService)

	zservice.AddDependService(
		etcdS.ZService,
		grpcClient,
		internal.Test_AssetUploadDownload().AddDependService(grpcClient),
		internal.Test_ConfigAssetUploadDownload().AddDependService(grpcClient),
	)

	zservice.Start().WaitStart()

}
