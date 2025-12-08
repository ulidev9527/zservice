package zecs

// 实体添加事件
type Event_EntityAdded struct {
	Entity int64
}

func (e *Event_EntityAdded) GetEventType() int32 { return Event_Type_EntityAdded }
