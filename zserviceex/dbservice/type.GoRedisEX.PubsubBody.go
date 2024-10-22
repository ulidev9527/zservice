package dbservice

// redis 消息订阅发布参数
type PubsubBody struct {
	S2S string `json:"s2s"` // 上下文
	Val []byte `json:"val"` // 内存数据
}
