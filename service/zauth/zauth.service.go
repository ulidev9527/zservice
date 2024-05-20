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
		internal.ZauthInitService = z
		internal.ZAuthInit()
	})

	internal.NsqService = nsqservice.NewNsqProducerService(&nsqservice.NsqProducerServiceConfig{
		Addr: zservice.Getenv("NSQD_ADDR"),
		OnStart: func(producer *nsq.Producer) {
			internal.Nsq = producer
			internal.InitNsq()
		},
	})

	internal.EtcdService = etcdservice.NewEtcdService(&etcdservice.EtcdServiceConfig{
		Addr: zservice.Getenv("ETCD_ADDR"),
		OnStart: func(etcd *clientv3.Client) {
			internal.Etcd = etcd
			internal.InitEtcd()
		},
	})

	internal.GrpcService = grpcservice.NewGrpcService(&grpcservice.GrpcServiceConfig{
		ListenAddr: zservice.Getenv("GRPC_LISTEN_ADDR"),
		EtcdServer: internal.Etcd,
		OnStart: func(grpc *grpc.Server) {
			internal.Grpc = grpc
			internal.InitGrpc()
		},
	})

	internal.GinService = ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenAddr: zservice.Getenv("GIN_LISTEN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			internal.Gin = engine
			internal.InitGin()
		},
	})

	zservice.AddDependService(
		internal.MysqlService.ZService, internal.RedisService.ZService, systemS,
		internal.NsqService.ZService, internal.EtcdService.ZService, internal.GrpcService.ZService,
		internal.GinService.ZService,
	)

	systemS.AddDependService(internal.MysqlService.ZService, internal.RedisService.ZService)

	internal.NsqService.AddDependService(systemS)

	internal.EtcdService.AddDependService(systemS)

	internal.GrpcService.AddDependService(systemS, internal.NsqService.ZService, internal.EtcdService.ZService)

	internal.GinService.AddDependService(internal.GrpcService.ZService)

	zservice.Start().WaitStart().WaitStop()
}
