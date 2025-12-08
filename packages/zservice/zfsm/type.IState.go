package zfsm

// IState 定义状态接口
type IState interface {
	GetType() int32  // 获取状态类型
	Enter()          // 状态机进入
	Update(dt int64) // 状态机更新
	Exit()           // 状态机离开
}
