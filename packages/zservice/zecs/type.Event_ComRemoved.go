package zecs

// 组件移除事件
type Event_ComRemoved struct {
	Entity    int64
	ComTypeID int32
}

func (e *Event_ComRemoved) GetEventType() int32 { return Event_Type_ComRemoved }
