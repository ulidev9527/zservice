package internal

import (
	"zservice/zservice/ex/redisservice"

	"gorm.io/gorm"
)

var (
	Mysql *gorm.DB
	Redis *redisservice.GoRedisEX
)

func InitMysql() {}
func InitRedis() {

}
