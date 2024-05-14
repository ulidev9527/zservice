package internal

import "gorm.io/gorm"

// 权限日志
type PermissionLogTable struct {
	gorm.Model
	UID     uint64 // 用户ID
	TraceID string // 链路ID
	PID     uint   // 权限ID
	Action  uint   // 动作
	Params  string // 参数记录,根据不同情况进行记录
}
