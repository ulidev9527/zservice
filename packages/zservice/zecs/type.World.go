package zecs

import (
	"slices"
	"sync"
	"time"
	"zserviceapps/packages/zservice"
)

// World 是 ECS 的世界，管理所有的 Archetype 和实体
type World struct {
	name           string
	archMgr        *ArchetypeMgr
	entityIDMax    int64 // 实体ID记录
	entityArcheMap map[int64]*Archetype
	lock           sync.RWMutex
	systems        []ISys
	cmdQueue       []command
	cmdLock        sync.Mutex
	lastUpdateTime time.Time // 新增字段，记录上次Update时间 ms
	entityPool     sync.Pool // 对象池: 复用 map[int32]ICom
	cmdPool        sync.Pool // 对象池: 复用 command

	isSysParallel bool            // 是否开启系统并行执行, 启用后实体越多数据表现越好
	sysParallelWG *sync.WaitGroup // 并行执行的工作组
	memoryBudget  int             // 内存预算，单位字节

	autoUpdateCD         time.Duration // 自动更新间隔时间, 默认 33，30帧每秒
	frameID              int64         // 帧更新数据
	isAutoUpdateRunning  bool          // 自动更新是否运行中
	isMarkStopAutoUpdate bool          // 是否标记停止自动更新

	eventBus *EventBus // 新增: 事件总线

	frameTime         time.Time // 当前帧时间
	frameTimestamp_ms int64     // 当前毫秒时间戳
}

type cmdType int

const (
	cmdAddEntity cmdType = iota
	cmdDestroyEntity
	cmdAddComponent
	cmdRemoveComponent
)

type command struct {
	typ     cmdType
	entity  int64
	comps   []ICom
	typeIDs []int32
}

const DefaultMemoryBudget = 16 * 1024 // 默认16KB内存预算

type WithNewWorldXXX func(*World)

func WithNewWorld_MemoryBudget(memoryBudget int) WithNewWorldXXX {
	return func(w *World) { w.memoryBudget = memoryBudget }
}
func WithNeewWorld_isSysParallel(isSysParallel bool) WithNewWorldXXX {
	return func(w *World) { w.isSysParallel = isSysParallel }
}
func WithNewWorld_AutoUpdateCD(d time.Duration) WithNewWorldXXX {
	return func(w *World) { w.autoUpdateCD = d }
}
func WithNewWorld_Name(name string) WithNewWorldXXX {
	return func(w *World) { w.name = name }
}

// 创建新 World (memoryBudget: 单个chunk的内存预算,单位字节,默认16KB)
func NewWorld(withNewWorldXXX ...WithNewWorldXXX) *World {

	w := &World{
		name:           "zecs",
		entityIDMax:    0,
		entityArcheMap: make(map[int64]*Archetype),
		systems:        make([]ISys, 0),
		cmdQueue:       make([]command, 0, 128),
		lastUpdateTime: time.Now(),
		entityPool: sync.Pool{
			New: func() interface{} {
				return make(map[int32]ICom, 4) // 预分配4个空间
			},
		},
		cmdPool: sync.Pool{
			New: func() interface{} {
				return &command{}
			},
		},
		sysParallelWG: &sync.WaitGroup{},
		eventBus:      NewEventBus(), // 新增
	}

	for _, with := range withNewWorldXXX {
		with(w)
	}

	if w.memoryBudget <= 0 {
		w.memoryBudget = DefaultMemoryBudget
	}

	w.archMgr = NewArchetypeMgr(w.memoryBudget)

	return w
}

// 开始自动更新
func (w *World) StartAutoUpdate() {

	if w.autoUpdateCD == 0 {
		w.autoUpdateCD = time.Millisecond * 33
	}

	if w.isAutoUpdateRunning {
		return
	}
	w.isAutoUpdateRunning = true
	zservice.Go(func() {

		for {
			if w.isMarkStopAutoUpdate {
				break
			}

			dt := w.Update() // 执行每帧逻辑
			if dt < w.autoUpdateCD {
				time.Sleep(w.autoUpdateCD - dt)
			} else {
				zservice.LogError(w.name, "frameDuration error", w.autoUpdateCD, dt)
			}
		}

		w.isAutoUpdateRunning = false

	})

}

// 停止自动更新
func (w *World) StopAutoUpdate() {
	if !w.isAutoUpdateRunning {
		return
	}
	w.isMarkStopAutoUpdate = true
	for {
		if !w.isAutoUpdateRunning {
			break
		}
	}
}

// 获取当前帧数
func (w *World) GetFrameID() int64 { return w.frameID }
func (w *World) ResetFrameInfo() {
	w.frameID = 1
}

// 创建新实体（ID生成在cmdLock中，避免死锁）
func (w *World) CreateEntity(comps ...ICom) int64 {
	w.cmdLock.Lock()
	defer w.cmdLock.Unlock()
	w.entityIDMax++
	entity := w.entityIDMax
	w.cmdQueue = append(w.cmdQueue, command{
		typ:    cmdAddEntity,
		entity: entity,
		comps:  comps,
	})
	return entity
}

// 销毁实体（改为异步，所有销毁都走命令队列）
func (w *World) RemoveEntity(entity int64) {
	w.cmdLock.Lock()
	w.cmdQueue = append(w.cmdQueue, command{
		typ:    cmdDestroyEntity,
		entity: entity,
	})
	w.cmdLock.Unlock()
}

// 添加组件（改为异步，所有组件添加都走命令队列）
func (w *World) AddCom(entity int64, comps ...ICom) {
	w.cmdLock.Lock()
	defer w.cmdLock.Unlock()
	// 组件去重（只保留最后一个同类型组件）
	compsMap := make(map[int32]ICom, len(comps))
	for _, com := range comps {
		compsMap[com.GetComType()] = com
	}
	uniqueComps := make([]ICom, 0, len(compsMap))
	for _, com := range compsMap {
		uniqueComps = append(uniqueComps, com)
	}
	w.cmdQueue = append(w.cmdQueue, command{
		typ:    cmdAddComponent,
		entity: entity,
		comps:  uniqueComps,
	})
}

// 移除组件（改为异步，所有组件移除都走命令队列）
func (w *World) RemoveCom(entity int64, typeIDs ...int32) {
	w.cmdLock.Lock()
	w.cmdQueue = append(w.cmdQueue, command{
		typ:     cmdRemoveComponent,
		entity:  entity,
		typeIDs: typeIDs,
	})
	w.cmdLock.Unlock()
}

// 执行命令队列（原地实现实体和组件的实际增删逻辑）
func (w *World) executeCmds() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.cmdLock.Lock()
	cmds := w.cmdQueue
	w.cmdQueue = make([]command, 0, 128)
	w.cmdLock.Unlock()

	for _, cmd := range cmds {
		switch cmd.typ {
		case cmdAddEntity:
			entity := cmd.entity
			ids := make([]int32, 0, len(cmd.comps))
			compsMap := make(map[int32]ICom, len(cmd.comps))
			for _, com := range cmd.comps {
				id := com.GetComType()
				ids = append(ids, id)
				compsMap[id] = com
			}
			arch := w.archMgr.GetOrCreateArchetype(ids)
			arch.AddEntity(entity, compsMap)
			w.entityArcheMap[int64(entity)] = arch
			w.SendEvent(&Event_EntityAdded{Entity: entity})
			for _, com := range compsMap {
				w.SendEvent(&Event_ComAdded{Entity: entity, ComTypeID: com.GetComType(), Com: com})
			}
		case cmdDestroyEntity:
			entity := cmd.entity
			arch := w.entityArcheMap[int64(entity)]
			if arch == nil {
				continue
			}
			arch.lock.Lock()
			for _, chunk := range arch.chunks {
				if chunk.HasEntity(entity) {
					chunk.lock.Lock()
					idx, ok := chunk.entityIndexMap[entity]
					if ok && idx < len(chunk.comps) {
						for typeID := range chunk.comps[idx] {
							w.SendEvent(&Event_ComRemoved{Entity: entity, ComTypeID: typeID})
						}
					}
					chunk.lock.Unlock()
					break
				}
			}
			arch.lock.Unlock()
			delete(w.entityArcheMap, int64(entity))
			arch.lock.Lock()
			for _, chunk := range arch.chunks {
				if chunk.HasEntity(entity) {
					chunk.lock.Lock()
					idx, ok := chunk.entityIndexMap[entity]
					if ok {
						lastIdx := chunk.count - 1
						if idx != lastIdx {
							chunk.entities[idx] = chunk.entities[lastIdx]
							chunk.comps[idx] = chunk.comps[lastIdx]
							chunk.entityIndexMap[chunk.entities[idx]] = idx
						}
						if lastIdx >= 0 {
							chunk.entities = chunk.entities[:lastIdx]
							chunk.comps = chunk.comps[:lastIdx]
						}
						delete(chunk.entityIndexMap, entity)
						chunk.count--
					}
					chunk.lock.Unlock()
					break
				}
			}
			arch.lock.Unlock()
			w.SendEvent(&Event_EntityRemoved{Entity: entity})
		case cmdAddComponent:
			entity := cmd.entity
			oldArch := w.entityArcheMap[int64(entity)]
			if oldArch == nil {
				continue
			}
			ids := make([]int32, len(oldArch.id))
			copy(ids, oldArch.id)
			chunk := oldArch.FindChunk(entity)
			if chunk == nil {
				continue
			}
			compsMap := chunk.GetEntityComps(entity)
			for _, com := range cmd.comps {
				id := com.GetComType()
				found := slices.Contains(ids, id)
				if !found {
					ids = append(ids, id)
				}
				compsMap[id] = com
			}
			newArch := w.archMgr.GetOrCreateArchetype(ids)
			// 检查新旧Archetype是否相同
			if newArch == oldArch {
				// 如果相同，只需要更新组件值，无需迁移实体
				for _, com := range cmd.comps {
					w.SendEvent(&Event_ComAdded{Entity: entity, ComTypeID: com.GetComType(), Com: com})
				}
				continue
			}
			// 不同Archetype时才需要迁移实体
			newArch.AddEntity(entity, compsMap)
			oldArch.RemoveEntity(entity)
			w.entityArcheMap[int64(entity)] = newArch
			for _, com := range cmd.comps {
				w.SendEvent(&Event_ComAdded{Entity: entity, ComTypeID: com.GetComType(), Com: com})
			}
		case cmdRemoveComponent:
			entity := cmd.entity
			oldArch := w.entityArcheMap[int64(entity)]
			if oldArch == nil {
				continue
			}
			chunk := oldArch.FindChunk(entity)
			if chunk == nil {
				continue
			}
			compsMap := chunk.GetEntityComps(entity)
			if compsMap == nil {
				continue
			}

			hasChanges := false
			removedIDs := make([]int32, 0, len(cmd.typeIDs))
			for _, id := range cmd.typeIDs {
				if _, exists := compsMap[id]; exists {
					delete(compsMap, id)
					hasChanges = true
					removedIDs = append(removedIDs, id)
				}
			}

			if !hasChanges {
				continue
			}

			ids := make([]int32, 0, len(compsMap))
			for id := range compsMap {
				ids = append(ids, id)
			}

			for _, id := range removedIDs {
				w.SendEvent(&Event_ComRemoved{Entity: entity, ComTypeID: id})
			}

			if len(ids) == 0 {
				oldArch.RemoveEntity(entity)
				delete(w.entityArcheMap, entity)
				w.SendEvent(&Event_EntityRemoved{Entity: entity})
				continue
			}

			newArch := w.archMgr.GetOrCreateArchetype(ids)
			newArch.AddEntity(entity, compsMap)
			oldArch.RemoveEntity(entity)
			w.entityArcheMap[int64(entity)] = newArch
		}
	}
}

// 遍历所有拥有指定组件类型的实体及其组件（只遍历包含全部 typeIDs 的 Archetype）
func (w *World) QueryCom(typeIDs []int32, fn func(entity int64, comps map[int32]ICom)) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	for _, arch := range w.archMgr.MatchArchetypes(typeIDs) {
		arch.Each(func(entity int64, comps map[int32]ICom) {
			// 只返回当前有效实体
			if _, ok := w.entityArcheMap[entity]; ok {
				fn(entity, comps)
			}
		})
	}
}

// 查询拥有指定组件类型的所有实体
func (w *World) Query(typeIDs []int32) []int64 {
	w.lock.RLock()
	defer w.lock.RUnlock()
	var result []int64
	for _, arch := range w.archMgr.MatchArchetypes(typeIDs) {
		for _, chunk := range arch.chunks {
			chunk.lock.RLock()
			result = append(result, chunk.entities[:chunk.count]...)
			chunk.lock.RUnlock()
		}
	}
	return result
}

// 查询拥有指定组件类型的所有实体及其组件集合
func (w *World) QueryComps(typeIDs []int32) []struct {
	Entity int64
	Comps  map[int32]ICom
} {
	w.lock.RLock()
	defer w.lock.RUnlock()
	result := make([]struct {
		Entity int64
		Comps  map[int32]ICom
	}, 0)
	for _, arch := range w.archMgr.MatchArchetypes(typeIDs) {
		for _, chunk := range arch.chunks {
			chunk.lock.RLock()
			for i := 0; i < chunk.count; i++ {
				entity := chunk.entities[i]
				// 只返回当前有效实体
				if _, ok := w.entityArcheMap[entity]; ok {
					result = append(result, struct {
						Entity int64
						Comps  map[int32]ICom
					}{
						Entity: entity,
						Comps:  chunk.comps[i],
					})
				}
			}
			chunk.lock.RUnlock()
		}
	}
	return result
}

// 判断实体是否拥有某组件
func (w *World) HasCom(entity int64, typeID int32) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	arch := w.entityArcheMap[int64(entity)]
	if arch == nil {
		return false
	}
	chunk := arch.FindChunk(entity)
	if chunk == nil {
		return false
	}
	comps := chunk.GetEntityComps(entity)
	_, ok := comps[typeID]
	return ok
}

// 获取指定实体的指定组件，若不存在返回nil
func (w *World) GetCom(entity int64, typeID int32) ICom {
	w.lock.RLock()
	defer w.lock.RUnlock()
	arch := w.entityArcheMap[int64(entity)]
	if arch == nil {
		return nil
	}
	chunk := arch.FindChunk(entity)
	if chunk == nil {
		return nil
	}
	comps := chunk.GetEntityComps(entity)
	if comps == nil {
		return nil
	}
	return comps[typeID]
}

// 注册系统
func (w *World) AddSys(sysList ...ISys) {

	// 按优先级插入
	for _, sys := range sysList {

		order := sys.GetOrder()
		idx := 0
		for i, s := range w.systems {
			if s.GetOrder() > order {
				idx = i
				break
			}
			idx = i + 1
		}

		// 在指定位置插入
		w.systems = append(w.systems, nil)
		copy(w.systems[idx+1:], w.systems[idx:])
		w.systems[idx] = sys
		sys.OnAdd(w)
	}
}

// 修改Update方法，无需参数
// @return 本次执行时长
func (w *World) Update() time.Duration {
	currentTime := time.Now()
	dt := currentTime.Sub(w.lastUpdateTime).Milliseconds() // 计算时间增量（秒）
	w.lastUpdateTime = currentTime
	w.UpdateDeltaTime(dt)
	return time.Since(currentTime)
}

// 新增：带dt参数的UpdateDeltaTime
func (w *World) UpdateDeltaTime(dt int64) {
	w.frameID++
	w.frameTime = time.Now()
	w.frameTimestamp_ms = w.frameTime.UnixMilli()
	w.SendNowEvent(&Event_FrameStart{frameID: w.frameID, dt: dt})

	// 第一阶段：系统更新（只读）
	w.lock.RLock()

	// 创建并发工作组

	for i := 0; i < len(w.systems); {
		if w.isSysParallel {
			// 查找连续的可并发执行的系统
			j := i
			for ; j < len(w.systems) && w.systems[j].IsParallel() && w.systems[j].GetOrder() == w.systems[i].GetOrder(); j++ {
			}

			// 如果找到多个可并发执行的系统
			if j > i+1 {
				for k := i; k < j; k++ {
					w.sysParallelWG.Add(1)
					sys := w.systems[k]
					zservice.Go(func() {
						defer w.sysParallelWG.Done()
						sys.Update(w, dt)
					})
				}
				w.sysParallelWG.Wait()
				i = j
				continue
			}
		}

		// 单个系统直接执行
		w.systems[i].Update(w, dt)
		i++
	}

	w.lock.RUnlock()

	// 第二阶段：执行增删命令
	w.executeCmds()

	// 第三阶段：统一触发所有事件
	w.eventBus.Flush()

	w.SendNowEvent(&Event_FrameEnd{frameID: w.frameID, dt: dt})
}

// 立即发送事件
func (w *World) SendNowEvent(event IEvent) { w.eventBus.SendNow(event) }

// 发送事件
func (w *World) SendEvent(event IEvent) { w.eventBus.Send(event) }

// 注册事件监听
func (w *World) OnEvent(eventType int32, fn EventListener) { w.eventBus.On(eventType, fn) }

// 获取事件总线（如需直接操作）
func (w *World) EventBus() *EventBus { return w.eventBus }

// 重置ID为 100
func (w *World) Reset100ID() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.entityIDMax = 1000
}

// 获取当前时间戳
func (w *World) NowTime() time.Time {
	return w.frameTime
}

// 获取当前年月
func (w *World) NowYearAndMonth() string {
	return w.frameTime.Format("2006-01")
}

// 当前时间戳 毫秒
func (w *World) NowTimesteamp_MS() int64 {
	return w.frameTimestamp_ms
}
