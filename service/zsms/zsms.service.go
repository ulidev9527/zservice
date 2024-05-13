package main

import (
	_ "embed"
	"zservice/service/zconfig/zconfig"
	"zservice/service/zsms/internal"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/redisservice"

	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func init() {

	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zsms",
		Version: "0.1.0",
	})

	if zservice.GetenvBool("USE_REMOTE_ENV") {
		e := zconfig.LoadRemoteEnv(zservice.Getenv("REMOTE_ENV_ADDR"), zservice.Getenv("REMOTE_ENV_AUTH"))
		if e != nil {
			zservice.LogPanic(e)
		}
	}

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
		OnStart: func(db *redis.Client) {
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
		Addr:       zservice.Getenv("GRPC_ADDR"),
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
