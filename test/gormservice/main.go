package main

import (
	"database/sql"
	"time"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"

	"gorm.io/gorm"
)

func init() {

	zservice.Init("gormservice_test", "1.0.0")
}

type TimeTable struct {
	gorm.Model
	T1 sql.NullTime
}

func main() {
	zservice.AddDependService(gormservice.NewGormMysqlService(&gormservice.GormMysqlServiceConfig{
		DBName: zservice.Getenv("MYSQL_DBNAME"),
		Addr:   zservice.Getenv("MYSQL_ADDR"),
		User:   zservice.Getenv("MYSQL_USER"),
		Pass:   zservice.Getenv("MYSQL_PASS"),
		Debug:  zservice.GetenvBool("MYSQL_DEBUG"),
		OnStart: func(s *gormservice.GormMysqlService) {
			s.Mysql.AutoMigrate(TimeTable{})

			// s.Mysql.Create(&TimeTable{
			// 	T1: sql.NullTime{Time: time.Now(), Valid: true},
			// })

			tabs := &TimeTable{}
			if e := s.Mysql.Find(tabs, "t1 < ?", time.Now()).Error; e != nil {
				s.LogError(e)
			} else {
				s.LogInfo(tabs)
			}

		},
	}).ZService).Start()
}
