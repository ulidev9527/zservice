package dbservice

import (
	"fmt"
	"zservice/zservice"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GOrmEXConfig struct {
	DBName string
	Addr   string // IP:PORT
	User   string
	Passwd string
	Debug  bool
}

func (conf *GOrmEXConfig) GetUri(isUsePasswd ...bool) string {

	if isUsePasswd != nil && isUsePasswd[0] {
		return fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", conf.User, conf.Passwd, conf.Addr, conf.DBName)
	} else {
		return fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", conf.User, "******", conf.Addr, conf.DBName)
	}
}

type GormMysqlService struct {
	*zservice.ZService
	Mysql *gorm.DB
}
type GormMysqlServiceConfig struct {
	Name    string         // 服务名称
	DBName  string         // 数据库名称
	Addr    string         // 数据库地址
	User    string         // 数据库用户名
	Pass    string         // 数据库密码
	OnStart func(*gorm.DB) // 启动的回调
}

func NewGormMysqlService(c *GormMysqlServiceConfig) *GormMysqlService {
	if c == nil {
		zservice.LogPanic("GormMysqlServiceConfig is nil")
		return nil
	}
	name := "GormMysqlService"
	if c.Name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	gs := &GormMysqlService{}
	zs := zservice.NewService(name, func(s *zservice.ZService) {
		db, e := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", c.User, c.Pass, c.Addr, c.DBName)), &gorm.Config{})
		if e != nil {
			zservice.LogPanic(e)
		}
		gs.Mysql = db

		if c.OnStart != nil {
			c.OnStart(gs.Mysql)
		}

		s.StartDone()
	})

	gs.ZService = zs
	return gs
}
