package zecs

// 帧开始事件
type Event_FrameStart struct {
	IEvent
	frameID int64
	dt      int64
}

func (ev *Event_FrameStart) GetEventType() int32 { return Event_Type_FrameStart }

func (ev *Event_FrameStart) GetFrameID() int64     { return ev.frameID }
func (ev *Event_FrameStart) GetDelayedTime() int64 { return ev.dt }
