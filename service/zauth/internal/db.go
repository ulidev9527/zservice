package internal

import (
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zhelper"

	"gorm.io/gorm"
)

var (
	MysqlService *gormservice.GormMysqlService
	Mysql        *gorm.DB
	RedisService *redisservice.RedisService
	Redis        *redisservice.GoRedisEX
	dbhelper     *zhelper.DBHelper
)

func InitMysql() {

	Mysql.AutoMigrate(AssetTable{})
	Mysql.AutoMigrate(OrgTable{})
	Mysql.AutoMigrate(PermissionBindTable{})
	Mysql.AutoMigrate(PermissionTable{})
	Mysql.AutoMigrate(ServiceKVTable{})
	Mysql.AutoMigrate(SmsBanTable{})
	Mysql.AutoMigrate(UserTable{})
	Mysql.AutoMigrate(UserOrgBindTable{})

	dbhelper = zhelper.NewDBHelper(Redis, Mysql)
}
func InitRedis() {

}
