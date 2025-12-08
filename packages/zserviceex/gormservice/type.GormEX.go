package gormservice

import (
	"errors"

	"gorm.io/gorm"
)

type GormEX struct {
	*gorm.DB
}

func (g *GormEX) IsNotFoundErr(e error) bool {
	return errors.Is(e, gorm.ErrRecordNotFound)
}
