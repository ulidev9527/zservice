package internal

import (
	"zservice/zservice/ex/gormservice"

	"gorm.io/gorm"
)

var (
	MysqlService *gormservice.GormMysqlService
	Mysql        *gorm.DB
)

func InitMysql() {
	Mysql.AutoMigrate(&LogKVTable{})
}
