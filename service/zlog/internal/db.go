package internal

import (
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"

	"gorm.io/gorm"
)

var (
	MysqlService *gormservice.GormMysqlService
	Mysql        *gorm.DB
	RedisService *redisservice.RedisService
	Redis        *redisservice.GoRedisEX
)

func InitMysql() {
	Mysql.AutoMigrate(&LogKVTable{})
}
func InitRedis() {

}
