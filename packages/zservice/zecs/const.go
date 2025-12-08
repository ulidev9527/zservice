package zecs

// 内置事件类型常量
const (
	Event_Type_EntityAdded   int32 = 1 // 实体添加
	Event_Type_EntityRemoved int32 = 2 // 实体移除
	Event_Type_ComAdded      int32 = 3 // 组件添加
	Event_Type_ComRemoved    int32 = 4 // 组件移除
	Event_Type_FrameStart    int32 = 5 // 帧逻辑开始
	Event_Type_FrameEnd      int32 = 6 // 帧逻辑结束
)
