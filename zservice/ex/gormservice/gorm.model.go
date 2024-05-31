package gormservice

import (
	"gorm.io/gorm"
)

type TimeModel struct {
	CreatedAt uint64         `gorm:"autoCreateTime:milli"`
	UpdatedAt uint64         `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AllModel struct {
	TimeModel
	ID uint `gorm:"primaryKey,autoIncrement"`
}
