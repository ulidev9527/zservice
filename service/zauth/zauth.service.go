package main

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_ex"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"
)

func init() {
	zservice.Init("zauth", "1.0.0")
}

func main() {

	internal.MysqlService = gormservice.NewGormMysqlService(&gormservice.GormMysqlServiceConfig{
		DBName: zservice.Getenv("MYSQL_DBNAME"),
		Addr:   zservice.Getenv("MYSQL_ADDR"),
		User:   zservice.Getenv("MYSQL_USER"),
		Pass:   zservice.Getenv("MYSQL_PASS"),
		Debug:  zservice.GetenvBool("MYSQL_DEBUG"),
		OnStart: func(s *gormservice.GormMysqlService) {
			internal.Mysql = s.Mysql
			internal.InitMysql()
		},
	})
	internal.RedisService = redisservice.NewRedisService(&redisservice.RedisServiceConfig{
		Addr: zservice.Getenv("REDIS_ADDR"),
		Pass: zservice.Getenv("REDIS_PASS"),
		OnStart: func(db *redisservice.GoRedisEX) {
			internal.Redis = db
			internal.InitRedis()
		},
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
	ServiceRegist.AddDependService(
		internal.MysqlService.ZService,
		internal.RedisService.ZService,
	)

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
