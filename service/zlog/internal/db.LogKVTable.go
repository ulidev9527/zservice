package internal

import "zservice/zservice/ex/gormservice"

// 键值对存储
type LogKVTable struct {
	gormservice.IDModel
	UID      uint32 // 用户ID
	Key      string // 键
	Value    string // 值
	SaveTime int64  // 保存时间
}
