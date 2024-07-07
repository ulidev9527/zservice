package internal

import (
	"zservice/zserviceex/dbservice"
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
