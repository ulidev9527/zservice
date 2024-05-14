package main

import (
	"zservice/service/zauth/internal"
	"zservice/service/zconfig/zconfig"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/ex/redisservice"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zauth",
		Version: "1.0.0",
	})
	e := zconfig.LoadRemoteEnv(zservice.Getenv("REMOTE_ENV_ADDR"), zservice.Getenv("REMOTE_ENV_AUTH"))
	if e != nil {
		zservice.LogPanic(e)
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

	nsqS := nsqservice.NewNsqProducerService(&nsqservice.NsqProducerServiceConfig{
		Addr: zservice.Getenv("NSQD_ADDR"),
		OnStart: func(producer *nsq.Producer) {
			internal.Nsq = producer
			internal.InitNsq()
		},
	})

	etcdS := etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{
		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			internal.Etcd = etcd
			internal.InitEtcd()
		},
	})

	grpcS := grpcservice.NewGrpcService(&grpcservice.GrpcServiceConfig{
		Addr:       zservice.Getenv("GRPC_ADDR"),
		EtcdServer: etcdS.Etcd,
		OnStart: func(grpc *grpc.Server) {
			internal.Grpc = grpc
			internal.InitGrpc()
		},
	})

	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		Addr: zservice.Getenv("GIN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			internal.Gin = engine
			internal.InitGin()
		},
	})

	zservice.AddDependService(mysqlS.ZService)
	zservice.AddDependService(redisS.ZService)
	zservice.AddDependService(ginS.ZService)
	zservice.AddDependService(etcdS.ZService)
	zservice.AddDependService(grpcS.ZService)
	zservice.AddDependService(nsqS.ZService)

	grpcS.AddDependService(mysqlS.ZService)
	grpcS.AddDependService(redisS.ZService)
	grpcS.AddDependService(nsqS.ZService)

	ginS.AddDependService(grpcS.ZService)
	ginS.AddDependService(etcdS.ZService)
	ginS.AddDependService(mysqlS.ZService)
	ginS.AddDependService(redisS.ZService)
	ginS.AddDependService(nsqS.ZService)

	zservice.Start()
	zservice.WaitStart()
	zservice.WaitStop()
}
