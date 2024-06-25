package internal

import "zservice/zservice/ex/gormservice"

// 键值对存储
type LogKVTable struct {
	gormservice.Model
	TraceID  string // 链路ID
	Key      string // 键
	Value    string // 值
	SaveTime int64  // 保存时间, 毫秒
	Service  string // 服务
}
