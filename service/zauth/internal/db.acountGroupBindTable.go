package internal

import "gorm.io/gorm"

// 账号组绑定
type AccountGroupBindTable struct {
	gorm.Model

	GID   uint   // 组ID
	UID   uint64 // 用户ID
	State uint   `gorm:"default:1"` // 状态 0禁用 1开启
}
