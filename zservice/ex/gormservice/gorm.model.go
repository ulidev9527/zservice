package gormservice

import (
	"gorm.io/gorm"
)

type Model struct {
	CreatedAt uint64         `gorm:"autoCreateTime:milli"`
	UpdatedAt uint64         `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `gorm:"primaryKey,autoIncrement"`
}
