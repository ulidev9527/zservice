package zservice

import "time"

// 心跳，在线保持 pro 版本
type Heartbeat struct {
	keepalive *Keepalive

	counter  int64 // 心跳计数
	delayed  int32 // 单次心跳延时
	tickTime int64 // 发送心跳的时间戳

	loseCount int64 // 丢失的心跳次数
}

func NewHeartbeat(ctx *Context, onTick func(ctx *Context)) *Heartbeat {
	h := &Heartbeat{}

	h.keepalive = NewKeepalive(ctx, func(ctx *Context) {
		h.counter++
		h.tickTime = time.Now().UnixMilli()
		onTick(ctx)
	})

	return h
}

// 设置心跳频率
func (h *Heartbeat) SetHeartbeatCD(cd int64) {
	h.keepalive.SetKeepaliveCD(cd)
}

// 更新心跳
func (h *Heartbeat) UpdateHeartbeat(dt int64) {
	h.keepalive.UpdateKeepalive(dt)
}

// 获取心跳次数
func (h *Heartbeat) GetHeartbeatCounter() int64 { return h.counter }

// 获取心跳延迟
func (h *Heartbeat) GetHeartbeatDelayed() int32 { return h.delayed }

// 回复心跳
// @return 返回丢失心跳的次数
func (h *Heartbeat) ReplyHeartbeat(count int64) int64 {
	if count != h.counter {
		h.loseCount++
		h.delayed = 999
		return h.loseCount
	}
	h.loseCount = 0
	dt := (time.Now().UnixMilli() - h.tickTime) / 2
	if dt > 999 {
		h.delayed = 999
	} else {
		h.delayed = int32(dt)
	}

	return 0
}

// 重置
func (h *Heartbeat) ResetHeartbeat() {

	h.loseCount = 0
	h.delayed = 0
	h.counter = 0
	h.keepalive.ResetKeepalive()
}
