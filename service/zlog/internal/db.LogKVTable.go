package internal

import (
	"time"

	"gorm.io/gorm"
)

// 键值对存储
type LogKVTable struct {
	gorm.Model
	TraceID  string    // 链路ID
	Key      string    // 键
	Value    string    // 值
	SaveTime time.Time // 保存时间, 毫秒
	Service  string    // 服务
}
