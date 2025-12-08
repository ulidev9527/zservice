package zecs

// 实体移除事件
type Event_EntityRemoved struct {
	Entity int64
}

func (e *Event_EntityRemoved) GetEventType() int32 { return Event_Type_EntityRemoved }
