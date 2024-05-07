package main

import (
	_ "embed"
	"os"
	"zservice/internal/dbservice"
	"zservice/service/zconfig/internal"
	"zservice/zservice"

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

	zservice.AddDependService(mysqlS.ZService)
	zservice.AddDependService(redisS.ZService)

	zservice.Start()
	zservice.WaitStart()
	zservice.WaitStop()
}
