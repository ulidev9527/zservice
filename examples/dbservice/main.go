package main

import (
	"database/sql"
	"time"

	"github.com/ulidev9527/zservice/zservice"
	"github.com/ulidev9527/zservice/zserviceex/dbservice"
	"gorm.io/gorm"
)

type TimeTestTable struct {
	gorm.Model

	// T1 DeletedAt
	// T1 zservice.TimeSQ
	T1 zservice.Time
	T2 sql.NullTime
}

func main() {

	zservice.Init(zservice.ZserviceOption{
		Name:    "gormservice.test",
		Version: "1.0.0",
	})

	zservice.AddDependService(
		dbservice.NewDBService(dbservice.DBServiceOption{
			DBType:      zservice.Getenv("DBSERVICE_GORM_TYPE"),
			DBName:      zservice.Getenv("DBSERVICE_GORM_NAME"),
			DBHost:      zservice.Getenv("DBSERVICE_GORM_HOST"),
			DBPort:      zservice.GetenvInt("DBSERVICE_GORM_PORT"),
			DBUser:      zservice.Getenv("DBSERVICE_GORM_USER"),
			DBPass:      zservice.Getenv("DBSERVICE_GORM_PASS"),
			RedisAddr:   zservice.Getenv("DBSERVICE_REDIS_ADDR"),
			RedisPass:   zservice.Getenv("DBSERVICE_REDIS_PASS"),
			RedisPrefix: zservice.Getenv("DBSERVICE_REDIS_PREFIX"),
			Debug:       zservice.GetenvBool("DBSERVICE_DEBUG"),
			OnStart: func(s *dbservice.DBService) {
				s.Gorm.AutoMigrate(TimeTestTable{})

				zservice.TestAction("insert", func() {

					tab := &TimeTestTable{
						T1: zservice.TimeNull(),
						// T1: DeletedAt{Time: time.Time{}},
						// T1: zservice.TimeSQ{Time: time.Time{}},
						T2: sql.NullTime{Time: time.Time{}},
					}

					if e := s.Gorm.Save(tab).Error; e != nil {
						s.LogError(e)
					} else {
						s.LogInfo(tab)
					}

				})

				zservice.TestAction("select", func() {
					tab := &TimeTestTable{}

					if e := s.Gorm.Order("created_at DESC").First(tab).Error; e != nil {
						s.LogError(e)
					} else {
						s.LogInfo(tab.T1.UnixMilli(), tab.T2.Time.UnixMilli())
					}
				})

			},
		}).ZService,
	).Start().WaitStop()
}
