package internal

import (
	"zservice/zservice/service/dbservice"
)

var (
	DBService *dbservice.DBService
	Gorm      *dbservice.GormEX
	Redis     *dbservice.GoRedisEX
)

func InitDB(s *dbservice.DBService) {

	DBService = s
	Gorm = DBService.Gorm
	Redis = DBService.Redis

	Gorm.AutoMigrate(
		LogKVTable{},
	)

}
