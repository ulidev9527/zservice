package zecs

import (
	"sync"
)

// 事件监听器
type EventListener func(event IEvent)

// 事件总线
type EventBus struct {
	listeners  map[int32][]EventListener
	lock       sync.RWMutex
	eventQueue []IEvent // 新增事件队列
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners:  make(map[int32][]EventListener),
		eventQueue: make([]IEvent, 0, 32),
	}
}

// Emit 立即触发
func (bus *EventBus) SendNow(event IEvent) {
	bus.lock.RLock()
	listeners := bus.listeners[event.GetEventType()]
	bus.lock.RUnlock()
	for _, fn := range listeners {
		fn(event)
	}
}

// Send 只入队
func (bus *EventBus) Send(event IEvent) {
	bus.lock.Lock()
	bus.eventQueue = append(bus.eventQueue, event)
	bus.lock.Unlock()
}

// Flush 统一触发所有事件
func (bus *EventBus) Flush() {
	bus.lock.Lock()
	queue := bus.eventQueue
	bus.eventQueue = make([]IEvent, 0, 32)
	bus.lock.Unlock()
	for _, event := range queue {
		bus.SendNow(event)
	}
}

func (bus *EventBus) On(eventType int32, fn EventListener) {
	bus.lock.Lock()
	bus.listeners[eventType] = append(bus.listeners[eventType], fn)
	bus.lock.Unlock()
}
