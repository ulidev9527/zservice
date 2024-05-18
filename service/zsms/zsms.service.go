package main

import (
	"zservice/service/zsms/internal"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func init() {
	zservice.Init("zsms", "1.0.0")
}

func main() {

	mysqlS := gormservice.NewGormMysqlService(&gormservice.GormMysqlServiceConfig{
		DBName: zservice.Getenv("MYSQL_DBNAME"),
		Addr:   zservice.Getenv("MYSQL_ADDR"),
		User:   zservice.Getenv("MYSQL_USER"),
		Pass:   zservice.Getenv("MYSQL_PASS"),
		OnStart: func(db *gorm.DB) {
			internal.Mysql = db
			internal.InitMysql()
		},
	})
	redisS := redisservice.NewRedisService(&redisservice.RedisServiceConfig{
		Addr: zservice.Getenv("REDIS_ADDR"),
		Pass: zservice.Getenv("REDIS_PASS"),
		OnStart: func(db *redisservice.GoRedisEX) {
			internal.Redis = db
			internal.InitRedis()
		},
	})

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{

		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			// do something
		},
	})

	grpcS := grpcservice.NewGrpcService(&grpcservice.GrpcServiceConfig{
		ListenAddr: zservice.Getenv("GRPC_ADDR"),
		EtcdServer: etcdS.Etcd,
		OnStart: func(grpc *grpc.Server) {
			zsms_pb.RegisterZsmsServer(grpc, internal.NewZsmsServer())
		},
	})

	zservice.AddDependService(mysqlS.ZService)
	zservice.AddDependService(redisS.ZService)
	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grpcS.ZService)

	grpcS.AddDependService(etcdS.ZService)
	grpcS.AddDependService(mysqlS.ZService)
	grpcS.AddDependService(redisS.ZService)

	zservice.Start()
	zservice.WaitStart()
	zservice.WaitStop()
}