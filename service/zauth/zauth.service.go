package main

import (
	"zservice/service/zauth/internal"
	"zservice/zservice"
	"zservice/zservice/ex/etcdservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/ex/nsqservice"
	"zservice/zservice/ex/redisservice"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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
		OnStart: func(db *gorm.DB) {
			internal.Mysql = db
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

	systemS := zservice.NewService("system", func(z *zservice.ZService) {
		internal.SystemService = z
		internal.SystemInit()
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
		ListenAddr: zservice.Getenv("GRPC_LISTEN_ADDR"),
		EtcdServer: etcdS.Etcd,
		OnStart: func(grpc *grpc.Server) {
			internal.Grpc = grpc
			internal.InitGrpc()
		},
	})

	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenAddr: zservice.Getenv("GIN_LISTEN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			internal.Gin = engine
			internal.InitGin()
		},
	})

	zservice.AddDependService(
		internal.MysqlService.ZService, internal.RedisService.ZService, systemS,
		nsqS.ZService, etcdS.ZService, grpcS.ZService,
		ginS.ZService,
	)

	systemS.AddDependService(internal.MysqlService.ZService, internal.RedisService.ZService)

	nsqS.AddDependService(systemS)

	etcdS.AddDependService(systemS)

	grpcS.AddDependService(systemS, nsqS.ZService, etcdS.ZService)

	ginS.AddDependService(grpcS.ZService)

	zservice.Start().WaitStart().WaitStop()
}
