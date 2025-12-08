package zservice

type StateMgr struct {
	state int32 // 状态

	stateChangeCallback func(ctx *Context, state int32) // 状态变化回调
	stateUpdateCallback func(state int32, dt int64)     // 状态更新回调
}

// 新建状态管理器
func NewStateMgr(stateChangeCallback func(ctx *Context, state int32), stateUpdateCallback func(state int32, dt int64)) *StateMgr {
	return &StateMgr{
		stateChangeCallback: stateChangeCallback,
		stateUpdateCallback: stateUpdateCallback,
	}
}

// 设置状态
// @return 是否改变
func (mgr *StateMgr) SetState(ctx *Context, state int32) bool {

	if mgr.state == state {
		return false
	}

	mgr.state = state

	// 状态变化回调
	mgr.stateChangeCallback(ctx, state)

	return true
}

// 获取状态
func (mgr *StateMgr) GetState() int32 {
	return mgr.state
}

// 状态更新
func (mgr *StateMgr) UpdateState(dt int64) {
	mgr.stateUpdateCallback(mgr.state, dt)
}
