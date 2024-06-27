package gormservice

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"zservice/zservice"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormMysqlService struct {
	*zservice.ZService
	Mysql *gorm.DB
}
type GormMysqlServiceConfig struct {
	DBName  string                  // 数据库名称
	Addr    string                  // 数据库地址
	User    string                  // 数据库用户名
	Pass    string                  // 数据库密码
	Debug   bool                    // 是否开启调试
	OnStart func(*GormMysqlService) // 启动的回调
}

func NewGormMysqlService(c *GormMysqlServiceConfig) *GormMysqlService {
	if c == nil {
		zservice.LogPanic("GormMysqlServiceConfig is nil")
		return nil
	}
	name := fmt.Sprint("GormMysqlService-", c.Addr, "-", c.DBName)

	gs := &GormMysqlService{}
	zs := zservice.NewService(name, func(s *zservice.ZService) {
		db, e := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true", c.User, c.Pass, c.Addr, c.DBName)), &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			}),
		})
		if e != nil {
			zservice.LogPanic(e)
		}

		if c.Debug {
			db = db.Debug()
		}

		gs.Mysql = db

		if c.OnStart != nil {
			c.OnStart(gs)
		}

		s.StartDone()
	})

	gs.ZService = zs
	return gs
}

// 是否是记录未找到
func IsNotFound(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}
