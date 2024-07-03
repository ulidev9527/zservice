package main

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_ex"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/service/dbservice"
	"zservice/zservice/service/etcdservice"
	"zservice/zservice/service/ginservice"
	"zservice/zservice/service/grpcservice"
	"zservice/zservice/zglobal"
)

func init() {
	zservice.Init("zauth", "1.0.0")
}

func main() {

	internal.DBService = dbservice.NewDBService(dbservice.DBServiceOption{
		GORMType:    zservice.Getenv("DBSERVICE_GORM_TYPE"),
		GORMName:    zservice.Getenv("DBSERVICE_GORM_NAME"),
		GORMAddr:    zservice.Getenv("DBSERVICE_GORM_ADDR"),
		GORMUser:    zservice.Getenv("DBSERVICE_GORM_USER"),
		GORMPass:    zservice.Getenv("DBSERVICE_GORM_PASS"),
		RedisAddr:   zservice.Getenv("DBSERVICE_REDIS_ADDR"),
		RedisPass:   zservice.Getenv("DBSERVICE_REDIS_PASS"),
		RedisPrefix: zservice.Getenv("DBSERVICE_REDIS_PREFIX"),
		Debug:       zservice.GetenvBool("DBSERVICE_DEBUG"),
		OnStart:     internal.InitDB,
	})

	ServiceRegist := zservice.NewService("ServiceRegist", func(z *zservice.ZService) {
		ctx := zservice.NewContext()
		zauth_ex.ServiceInfo.Regist(ctx, &zauth_pb.ServiceRegist_REQ{
			InitPermissions: []*zauth_pb.PermissionInfo{
				{Action: "post", Path: "/login", State: zglobal.E_PermissionState_AllowAll},
			},
		}, true)

		z.StartDone()
	})
	ServiceRegist.AddDependService(internal.DBService.ZService)

	internal.EtcdService = etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{
		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(s *etcdservice.EtcdService) {
			internal.Etcd = s.EtcdClient
			internal.InitEtcd()
		},
	})
	internal.EtcdService.AddDependService(ServiceRegist)

	internal.GrpcService = grpcservice.NewGrpcService(&grpcservice.GrpcServiceConfig{
		ListenPort: zservice.Getenv("grpc_listen_port"),
		EtcdClient: internal.EtcdService.EtcdClient,
		OnStart: func(s *grpcservice.GrpcService) {
			internal.Grpc = s.GrpcServer
			internal.InitGrpc()
		},
	})
	internal.GrpcService.AddDependService(internal.EtcdService.ZService)

	internal.GinService = ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenPort: zservice.Getenv("gin_listen_port"),
		OnStart: func(s *ginservice.GinService) {
			s.Engine.Use(zauth_ex.GinCheckAuthMiddleware(internal.GinService.ZService, true))
			internal.Gin = s.Engine
			internal.InitGin()
		},
	})
	internal.GinService.AddDependService(internal.GrpcService.ZService)

	readyS := zservice.NewService("ready", func(z *zservice.ZService) {
		internal.InitZZZZ()
		z.StartDone()
	})
	readyS.AddDependService(internal.GinService.ZService)

	zservice.AddDependService(readyS)

	zservice.Start().WaitStart().WaitStop()
}
