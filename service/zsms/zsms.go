package main

import (
	_ "embed"
	"zservice/internal/etcdservice"
	"zservice/internal/gormservice"
	"zservice/internal/grpcservice"
	"zservice/internal/redisservice"
	"zservice/service/zsms/internal"
	"zservice/service/zsms/zsms_pb"
	"zservice/zservice"

	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

//go:embed version
var Version string

func init() {

	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zsms",
		Version: Version,
	})

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