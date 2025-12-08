package zservice

// 保持在线
type Keepalive struct {
	loopCD int64              // 设置的更新频率
	cd     int64              // 在线保持 CD (ms), 保持 10s， 每 3s 更新一次
	onLoop func(ctx *Context) // 循环回调
}

func NewKeepalive(ctx *Context, onLoop func(ctx *Context)) *Keepalive {
	return &Keepalive{
		loopCD: 5000,
		cd:     0,
		onLoop: onLoop,
	}
}

// 设置更新频率
func (k *Keepalive) SetKeepaliveCD(cd int64) {
	k.loopCD = cd
}

// 更新在线信息
func (k *Keepalive) UpdateKeepalive(dt int64) {
	k.cd -= dt
	if k.cd > 0 {
		return
	}
	k.CallKeepliveLoop(GetMainCtx())
}

// 触发在线保持回调
func (k *Keepalive) CallKeepliveLoop(ctx *Context) {
	k.cd = k.loopCD
	k.onLoop(ctx)
}

// 重置
func (k *Keepalive) ResetKeepalive() {
	k.cd = 0
}
