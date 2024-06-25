package gormservice

import (
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey,autoIncrement"`
	CreatedAt uint64         `gorm:"autoCreateTime:milli"`
	UpdatedAt uint64         `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
