package main

import (
	_ "embed"
	"zservice/internal/ginservice"
	"zservice/internal/gormservice"
	"zservice/internal/redisservice"
	"zservice/service/zconfig/internal"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zconfig",
		Version: "0.1.0",
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

	ginS.AddDependService(mysqlS.ZService)
	ginS.AddDependService(redisS.ZService)

	zservice.Start()
	zservice.WaitStart()
	zservice.WaitStop()
}
