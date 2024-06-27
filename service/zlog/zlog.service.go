package main

import (
	"zservice/service/zlog/internal"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
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
		OnStart: func(s *gormservice.GormMysqlService) {
			internal.Mysql = s.Mysql
			internal.InitMysql()
		},
	})

	zservice.AddDependService(
		internal.MysqlService.ZService,
	)

	zservice.Start().WaitStart()

	internal.NsqInit()

	zservice.WaitStop()
}
