package internal

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Mysql *gorm.DB
	Redis *redis.Client
)

func InitMysql() {

	Mysql.AutoMigrate(&AccountTable{})
	Mysql.AutoMigrate(&AccountGroupTable{})
	Mysql.AutoMigrate(&AccountGroupBindTable{})
	Mysql.AutoMigrate(&PermissionTable{})
	Mysql.AutoMigrate(&PermissionLogTable{})
	Mysql.AutoMigrate(&PermissionBindTable{})
}
func InitRedis() {

}
