package zecs

type ICom interface {
	GetComType() int32 // 获取组件类型，需要自行实现
}
