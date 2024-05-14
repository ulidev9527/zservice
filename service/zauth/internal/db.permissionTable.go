package internal

import "gorm.io/gorm"

// 权限
type PermissionTable struct {
	gorm.Model
	Name   string // 权限名称
	Action string // 权限动作 详情：(readme.md#permissiontableaction)
	State  uint   `gorm:"default:2"` // 默认权限状态 0 禁用 1 开启 2 继承，继承父级权限
	PID    uint   // 父级ID
}
