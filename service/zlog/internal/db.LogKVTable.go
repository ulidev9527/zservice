package internal

// 键值对存储
type LogKVTable struct {
	TraceID  string // 链路ID
	UID      uint32 // 用户ID
	Key      string // 键
	Value    string // 值
	SaveTime int64  // 保存时间, 毫秒
}
