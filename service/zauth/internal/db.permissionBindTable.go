package internal

import (
	"time"

	"gorm.io/gorm"
)

type PermissionBindTable struct {
	gorm.Model
	GID     uint       // 组ID
	PID     uint       // 权限ID
	Expires *time.Time // 过期
	State   uint       // 状态 0禁用 1开启
}
