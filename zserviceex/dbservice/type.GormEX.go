package dbservice

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

type GormEX struct {
	*gorm.DB
}

func NewGormEX(opt DBServiceOption) *GormEX {

	db := &gorm.DB{}

	switch opt.GORMType {
	default:
		_db, e := gorm.Open(
			mysql.Open(
				fmt.Sprintf("%s:%s@tcp(%v)/%s?parseTime=true",
					opt.GORMUser,
					opt.GORMPass,
					opt.GORMAddr,
					opt.GORMName),
			),
			&gorm.Config{
				Logger: logger.New(
					log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
						SlowThreshold:             200 * time.Millisecond,
						LogLevel:                  logger.Warn,
						IgnoreRecordNotFoundError: true,
						Colorful:                  true,
					},
				),
			},
		)

		if e != nil {
			zservice.LogPanic(e)
		}
		db = _db
	}

	if opt.Debug {
		db = db.Debug()
	}
	return &GormEX{
		DB: db,
	}
}

func (ex *GormEX) IsNotFoundErr(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}
