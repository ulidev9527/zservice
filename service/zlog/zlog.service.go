package main

import (
	"zservice/service/zlog/internal"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"
)

func init() {
	zservice.Init("zlog", "1.0.0")
}

func main() {

	internal.DBService = dbservice.NewDBService(dbservice.DBServiceOption{
		GORMType:    zservice.Getenv("DBSERVICE_GORM_TYPE"),
		GORMName:    zservice.Getenv("DBSERVICE_GORM_NAME"),
		GORMAddr:    zservice.Getenv("DBSERVICE_GORM_ADDR"),
		GORMUser:    zservice.Getenv("DBSERVICE_GORM_USER"),
		GORMPass:    zservice.Getenv("DBSERVICE_GORM_PASS"),
		RedisAddr:   zservice.Getenv("DBSERVICE_REDIS_ADDR"),
		RedisPass:   zservice.Getenv("DBSERVICE_REDIS_PASS"),
		RedisPrefix: zservice.Getenv("DBSERVICE_REDIS_PREFIX"),
		Debug:       zservice.GetenvBool("DBSERVICE_DEBUG"),
		OnStart:     internal.InitDB,
	})

	zservice.AddDependService(
		internal.DBService.ZService,
	)

	zservice.Start().WaitStart()

	internal.NsqInit()

	zservice.WaitStop()
}
