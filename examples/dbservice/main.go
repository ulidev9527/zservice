package main

import (
	"database/sql"
	"time"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/ulidev9527/zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

func init() {
	zservice.Init("gormservice_test", "1.0.0")
}

type TimeTestTable struct {
	gorm.Model

	// T1 DeletedAt
	// T1 zservice.TimeSQ
	T1 zservice.Time
	T2 sql.NullTime
}

func main() {
	zservice.AddDependService(
		dbservice.NewDBService(dbservice.DBServiceOption{
			GORMType:    zservice.Getenv("DBSERVICE_GORM_TYPE"),
			GORMName:    zservice.Getenv("DBSERVICE_GORM_NAME"),
			GORMAddr:    zservice.Getenv("DBSERVICE_GORM_ADDR"),
			GORMUser:    zservice.Getenv("DBSERVICE_GORM_USER"),
			GORMPass:    zservice.Getenv("DBSERVICE_GORM_PASS"),
			RedisAddr:   zservice.Getenv("DBSERVICE_REDIS_ADDR"),
			RedisPass:   zservice.Getenv("DBSERVICE_REDIS_PASS"),
			RedisPrefix: zservice.Getenv("DBSERVICE_REDIS_PREFIX"),
			Debug:       zservice.GetenvBool("DBSERVICE_DEBUG"),
			OnStart: func(s *dbservice.DBService) {
				s.Gorm.AutoMigrate(TimeTestTable{})

				zservice.TestAction("insert", func() {

					tab := &TimeTestTable{
						T1: zservice.TimeNull(),
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
