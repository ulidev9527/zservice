package main

import (
	"zservice/service/zlog/internal"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	zservice.Init("zlog", "1.0.0")
}

func main() {

	internal.MysqlService = gormservice.NewGormMysqlService(&gormservice.GormMysqlServiceConfig{
		DBName: zservice.Getenv("MYSQL_DBNAME"),
		Addr:   zservice.Getenv("MYSQL_ADDR"),
		User:   zservice.Getenv("MYSQL_USER"),
		Pass:   zservice.Getenv("MYSQL_PASS"),
		Debug:  zservice.GetenvBool("MYSQL_DEBUG"),
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

	internal.GinService = ginservice.NewGinService(&ginservice.GinServiceConfig{
		ListenAddr: zservice.Getenv("GIN_LISTEN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			internal.Gin = engine
			internal.InitGin()
		},
	})

	internal.GinService.AddDependService(internal.MysqlService.ZService, internal.RedisService.ZService)

	zservice.AddDependService(
		internal.MysqlService.ZService,
		internal.RedisService.ZService,
		internal.GinService.ZService,
	)

	zservice.Start().WaitStart()

	internal.NsqInit()

	zservice.WaitStop()
}
