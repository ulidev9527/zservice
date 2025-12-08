package zecs

// 组件添加事件
type Event_ComAdded struct {
	Entity    int64
	ComTypeID int32
	Com       ICom
}

func (e *Event_ComAdded) GetEventType() int32 { return Event_Type_ComAdded }
