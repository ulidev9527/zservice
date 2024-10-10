package dbservice

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ulidev9527/zservice/zservice"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Added SQLite driver import
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormEX struct {
	*gorm.DB
}

func NewGormEX(opt DBServiceOption) *GormEX {

	db := &gorm.DB{}

	var dt gorm.Dialector

	switch opt.DBType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			opt.DBUser,
			opt.DBPass,
			opt.DBHost,
			opt.DBPort,
			opt.DBName,
		)
		if opt.DBParams != "" {
			dsn = dsn + "&" + opt.DBParams
		}
		dt = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			opt.DBHost,
			opt.DBUser,
			opt.DBPass,
			opt.DBName,
			opt.DBPort,
		)
		if opt.DBParams != "" {
			dsn = dsn + " " + opt.DBParams
		}
		dt = postgres.Open(dsn)

	case "sqlite":
		dt = sqlite.Open(opt.DBHost)

	default:
		zservice.LogPanic("not support gorm type", opt.DBType)
	}

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
		zservice.LogPanic(e)
	} else {
		db = _db
	}

	if opt.Debug {
		db = db.Debug()
	}

	if sdb, e := db.DB(); e != nil {
		zservice.LogPanic(e)
	} else {

		sdb.SetMaxIdleConns(opt.MaxIdleConns)

		sdb.SetMaxOpenConns(opt.MaxOpenConns)

		sdb.SetConnMaxLifetime(time.Duration(opt.ConnMaxLifetime) * time.Second)
	}

	return &GormEX{
		DB: db,
	}
}

func (ex *GormEX) IsNotFoundErr(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}
