package gormservice

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"zserviceapps/packages/zservice"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service struct {
	zservice *zservice.ZService
	GormEx   *GormEX // 数据库 gorm

}

type ServiceOption struct {
	DebugModel int // zserive.Debug_Model_XXX

	Type            string  // 数据库 类型 目前支持 mysql/postgres
	Name            string  // 数据库 名称
	Host            string  // 数据库 地址 填入地址才会启用 Gorm 功能
	Port            int     // 数据库 端口
	User            string  // 数据库 用户名
	Pass            string  // 数据库 密码
	Params          string  // 数据库 额外参数
	MaxIdleConns    int     // 最大空闲连接数 default: 10
	MaxOpenConns    int     // 最大连接数 default: 30
	ConnMaxLifetime float32 // 连接最大生命周期 default: 300s

	OnStart func(*Service)
}

// 同步初始化参数
func syncDefaultOption(opt ServiceOption) ServiceOption {
	if opt.MaxIdleConns == 0 {
		opt.MaxIdleConns = 10
	}
	if opt.MaxOpenConns == 0 {
		opt.MaxOpenConns = 30
	}
	if opt.ConnMaxLifetime == 0 {
		opt.ConnMaxLifetime = 3
	}
	return opt
}

// 创建对应数据库服务
func initDT(ser *Service, opt ServiceOption) gorm.Dialector {

	var dt gorm.Dialector
	switch opt.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			opt.User,
			opt.Pass,
			opt.Host,
			opt.Port,
			opt.Name,
		)
		if opt.Params != "" {
			dsn = dsn + "&" + opt.Params
		}
		dt = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			opt.Host,
			opt.User,
			opt.Pass,
			opt.Name,
			opt.Port,
		)
		if opt.Params != "" {
			dsn = dsn + " " + opt.Params
		}
		dt = postgres.Open(dsn)

	case "sqlite":
		dt = sqlite.Open(opt.Host)

	default:
		ser.zservice.LogError("not support gorm type", opt.Type)
		os.Exit(1)
	}
	return dt
}

func onStart(ser *Service, opt ServiceOption) *gorm.DB {
	dt := initDT(ser, opt)

	for {
		var db *gorm.DB
		// 连接数据库
		if _db, e := gorm.Open(dt, &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					LogLevel:                  logger.Warn,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
		}); e != nil {
			ser.zservice.LogError("has error, waiting 5s again", e)
			time.Sleep(time.Second * 5)
			continue
		} else {
			db = _db
		}
		if ser.zservice.IsDebug {
			db = db.Debug()
		}

		// 鬼知道干什么
		if sdb, e := db.DB(); e != nil {

			ser.zservice.LogError("has error, waiting 5s again", e)
			time.Sleep(time.Second * 5)
			continue

		} else {

			sdb.SetMaxIdleConns(opt.MaxIdleConns)

			sdb.SetMaxOpenConns(opt.MaxOpenConns)

			sdb.SetConnMaxLifetime(time.Duration(opt.ConnMaxLifetime) * time.Second)
		}

		return db

	}
}

func NewService(opt ServiceOption) *Service {
	opt = syncDefaultOption(opt)
	ser := &Service{}

	ser.zservice = zservice.NewService(zservice.ServiceOptions{
		Name: "gormservice-" + opt.Name,
		OnStart: func(s *zservice.ZService) {
			ser.GormEx = &GormEX{
				DB: onStart(ser, opt),
			}
			if opt.OnStart != nil {
				opt.OnStart(ser)
			}
		},
	})
	if opt.DebugModel != zservice.Debug_Model_Auto {
		ser.zservice.IsDebug = opt.DebugModel == zservice.Debug_Model_On
	}

	return ser
}
func IsNotFoundErr(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}

func (ex *Service) IsNotFoundErr(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}

func (ser *Service) GetZService() *zservice.ZService { return ser.zservice }
