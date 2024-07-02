package etcdservice

type EventBody struct {
	S2S string `json:"s2s"` // S2S 数据
	Val string `json:"val"` // 内存
}
