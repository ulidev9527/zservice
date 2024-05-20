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

	Mysql.AutoMigrate(&ZauthAccountTable{})
	Mysql.AutoMigrate(&ZauthOrgTable{})
	Mysql.AutoMigrate(&ZauthPermissionTable{})
	Mysql.AutoMigrate(&ZauthAccountOrgBindTable{})
	Mysql.AutoMigrate(&ZauthPermissionBindTable{})

	dbhelper = zhelper.NewDBHelper(Redis, Mysql)
}
func InitRedis() {

	// if e := Mysql.Raw(`
	// WITH RECURSIVE cte(id) AS (
	// 	SELECT g_id FROM account_group_bind_tables WHERE uid=?
	// 	UNION ALL SELECT
	// 	agt.g_id FROM cte JOIN account_group_tables agt ON cte.id = agt.id
	// ) SELECT DISTINCT id FROM cte WHERE id > 0;
	// `, 1001).Find(&[]struct{}{}).Error; e != nil {
	// 	zservice.LogError(e)
	// }
}
