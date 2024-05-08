package main

import (
	_ "embed"
	"os"
	"zservice/internal/ginservice"
	"zservice/internal/gormservice"
	"zservice/internal/redisservice"
	"zservice/service/zconfig/internal"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

//go:embed version
var Version string

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zconfig",
		Version: Version,
	})
}

func main() {

	mysqlS := gormservice.NewGormMysqlService(&gormservice.GormMysqlServiceConfig{
		DBName: os.Getenv("MYSQL_DBNAME"),
		Addr:   os.Getenv("MYSQL_ADDR"),
		User:   os.Getenv("MYSQL_USER"),
		Pass:   os.Getenv("MYSQL_PASS"),
		OnStart: func(db *gorm.DB) {
			internal.Mysql = db
			internal.InitMysql()
		},
	})
	redisS := redisservice.NewRedisService(&redisservice.RedisServiceConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		Pass: os.Getenv("REDIS_PASS"),
		OnStart: func(db *redis.Client) {
			internal.Redis = db
			internal.InitRedis()
		},
	})

	ginS := ginservice.NewGinService(&ginservice.GinServiceConfig{
		Addr: os.Getenv("GIN_ADDR"),
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
