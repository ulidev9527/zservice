package zfsm

type WatchStateChangeFN func(state IState)

// StateMachine 状态机基础结构
type StateMachine struct {
	lastState        IState // 上一个状态
	currentState     IState
	nextState        IState             // 下一个状态
	watchStateChange WatchStateChangeFN // 状态改变响应
}

// 获取当前状态
func (sm *StateMachine) GetCurrentStateType() int32 {
	if sm.currentState == nil {
		return 0
	}
	return sm.currentState.GetType()
}

// 状态改变
func (sm *StateMachine) SetState(state IState) {
	sm.nextState = state
}

// 尝试设置状态，如果已经有则不设置
func (sm *StateMachine) TrySetState(state IState) bool {
	if sm.nextState != nil {
		return false
	}
	sm.SetState(state)
	return true
}

// 是否有下一个状态
func (sm *StateMachine) HasNextState() bool {
	return sm.nextState != nil
}

// 获取上一个状态
func (sm *StateMachine) GetLastState() IState {
	return sm.lastState
}

// 配置状态改变响应
func (sm *StateMachine) WatchStateChange(fn WatchStateChangeFN) {
	sm.watchStateChange = fn
}

// 状态更新
func (sm *StateMachine) Update(dt int64) {

	if sm.nextState != nil { // 新状态先切换状态
		sm.lastState = sm.currentState
		if sm.currentState != nil {
			sm.currentState.Exit()
		}
		sm.currentState = nil

		if sm.nextState != nil {
			sm.currentState = sm.nextState
		}
		sm.currentState.Enter()
		if sm.watchStateChange != nil {
			sm.watchStateChange(sm.currentState)
		}
		sm.nextState = nil // 清空下一个状态
	}

	if sm.currentState != nil {
		sm.currentState.Update(dt)
	}

}
