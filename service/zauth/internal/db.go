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

	Mysql.AutoMigrate(&AssetTable{})
	Mysql.AutoMigrate(&OrgTable{})
	Mysql.AutoMigrate(&UserTable{})
	Mysql.AutoMigrate(&PermissionBindTable{})
	Mysql.AutoMigrate(&PermissionTable{})
	Mysql.AutoMigrate(&UserOrgBindTable{})
	Mysql.AutoMigrate(&SmsBanTable{})

	dbhelper = zhelper.NewDBHelper(Redis, Mysql)
}
func InitRedis() {

	// if e := Mysql.Raw(`
	// WITH RECURSIVE cte(id) AS (
	// 	SELECT g_id FROM user_group_bind_tables WHERE uid=?
	// 	UNION ALL SELECT
	// 	agt.g_id FROM cte JOIN user_group_tables agt ON cte.id = agt.id
	// ) SELECT DISTINCT id FROM cte WHERE id > 0;
	// `, 1001).Find(&[]struct{}{}).Error; e != nil {
	// 	zservice.LogError(e)
	// }
}
