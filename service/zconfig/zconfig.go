package main

import (
	_ "embed"
	"os"
	"zservice/internal/dbservice"
	"zservice/internal/httpservice"
	"zservice/service/zconfig/internal"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

//go:embed version
var Version string

func init() {
	zservice.LogDebug()
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "zconfig",
		Version: Version,
	})
}

func main() {

	mysqlS := dbservice.NewGormMysqlService(&dbservice.GormMysqlServiceConfig{
		DBName: os.Getenv("MYSQL_DBNAME"),
		Addr:   os.Getenv("MYSQL_ADDR"),
		User:   os.Getenv("MYSQL_USER"),
		Pass:   os.Getenv("MYSQL_PASS"),
		OnStart: func(db *gorm.DB) {
			internal.Mysql = db
		},
	})
	redisS := dbservice.NewRedisService(&dbservice.RedisServiceConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		Pass: os.Getenv("REDIS_PASS"),
		OnStart: func(db *redis.Client) {
			internal.Redis = db
		},
	})

	ginS := httpservice.NewGinService(&httpservice.GinServiceConfig{
		Addr: os.Getenv("GIN_ADDR"),
		OnStart: func(engine *gin.Engine) {
			internal.Gin = engine
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
