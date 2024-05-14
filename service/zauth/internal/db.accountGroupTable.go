package internal

import "gorm.io/gorm"

// 账号组
type AccountGroupTable struct {
	gorm.Model

	Name  string // 组名
	GID   uint   // 父级组ID
	State uint   `gorm:"default:1"` // 状态 0 禁用 1 开启

}
