package internal

import (
	"time"

	"gorm.io/gorm"
)

// 账号权限绑定表
type ZauthPermissionBindTable struct {
	gorm.Model
	PermissionBindID uint       // 权限绑定ID
	OtherID          uint       // 外部ID
	OtherIDType      uint       // 外部ID类型 0无效 1组织 2账号
	PermissionID     uint       // 权限ID
	PerentID         uint       // 父级权限绑定ID
	Expires          *time.Time // 过期时间
	State            uint       `gorm:"default:2"` // 状态 0禁止访问 1允许访问 2继承父级
}
