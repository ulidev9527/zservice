package zecs

// 帧结束事件
type Event_FrameEnd struct {
	IEvent
	frameID int64
	dt      int64
}

func (ev *Event_FrameEnd) GetEventType() int32 { return Event_Type_FrameEnd }

func (ev *Event_FrameEnd) GetFrameID() int64 { return ev.frameID }

func (ev *Event_FrameEnd) GetDelayedTime() int64 { return ev.dt }
