package zecs

import (
	"testing"
	"time"
)

// 示例组件
type Pos struct {
	X, Y int
}

func (p *Pos) GetComType() int32 { return 1 }

type Vel struct {
	VX, VY int
}

func (v *Vel) GetComType() int32 { return 2 }

// 示例系统
type MoveSys struct{ ISys }

func (ms *MoveSys) Update(world *World, dt int64) {
	world.QueryCom([]int32{1, 2}, func(entity int64, comps map[int32]ICom) {
		pos := comps[1].(*Pos)
		vel := comps[2].(*Vel)
		pos.X += vel.VX
		pos.Y += vel.VY
	})
}
func (ms *MoveSys) GetOrder() int32  { return 0 }
func (ms *MoveSys) IsParallel() bool { return true }

type HealthSys struct{ ISys }

func (hs *HealthSys) Update(world *World, dt int64) {
	world.QueryCom([]int32{3}, func(entity int64, comps map[int32]ICom) {
		h := comps[3].(*Health)
		h.HP++
	})
}
func (hs *HealthSys) GetOrder() int32  { return 1 }
func (hs *HealthSys) IsParallel() bool { return true }

type TagSys struct{ ISys }

func (ts *TagSys) Update(world *World, dt int64) {
	world.QueryCom([]int32{4}, func(entity int64, comps map[int32]ICom) {
		t := comps[4].(*Tag)
		t.T++
	})
}
func (ts *TagSys) GetOrder() int32  { return 1 }
func (ts *TagSys) IsParallel() bool { return true }

// 添加系统顺序测试
func TestECS_SystemOrder(t *testing.T) {
	world := NewWorld()
	ts := &TagSys{}
	hs := &HealthSys{}
	ms := &MoveSys{}

	// 乱序添加
	world.AddSys(ts)
	world.AddSys(ms)
	world.AddSys(hs)

	// 验证优先级顺序
	if len(world.systems) != 3 {
		t.Fatal("系统数量错误")
	}

	// MoveSys应该在最前面(order=0)
	if world.systems[0] != ms {
		t.Error("MoveSys优先级排序错误")
	}

	// HealthSys和TagSys的顺序应该和添加顺序一致(order=1)
	if world.systems[1] != ts || world.systems[2] != hs {
		t.Error("相同优先级系统顺序错误")
	}
}

func TestECS_Basic(t *testing.T) {
	t.Run("创建实体和组件", func(t *testing.T) {
		world := NewWorld() // 使用默认内存预算
		pos := &Pos{X: 0, Y: 0}
		vel := &Vel{VX: 1, VY: 2}
		e := world.CreateEntity(pos, vel)
		world.Update() // 使NewEntity生效

		if !world.HasCom(e, 1) {
			t.Error("无法找到Position组件")
		}
		if !world.HasCom(e, 2) {
			t.Error("无法找到Velocity组件")
		}

		comps := world.QueryComps([]int32{1, 2})
		if len(comps) != 1 {
			t.Fatal("查询组件数量错误")
		}
		p := comps[0].Comps[1].(*Pos)
		v := comps[0].Comps[2].(*Vel)
		if p.X != 0 || p.Y != 0 {
			t.Errorf("Position初始值错误: got %+v", p)
		}
		if v.VX != 1 || v.VY != 2 {
			t.Errorf("Velocity初始值错误: got %+v", v)
		}
	})

	t.Run("添加和更新组件", func(t *testing.T) {
		world := NewWorld() // 使用默认内存预算
		e := world.CreateEntity(&Pos{X: 0, Y: 0})
		world.Update() // 使NewEntity生效
		world.AddCom(e, &Vel{VX: 1, VY: 2})
		world.Update() // 使AddCom生效

		if !world.HasCom(e, 2) {
			t.Fatal("添加组件失败")
		}

		// 更新已存在的组件
		world.AddCom(e, &Vel{VX: 3, VY: 4})
		world.Update() // 使AddCom生效
		comps := world.QueryComps([]int32{2})
		if len(comps) != 1 {
			t.Fatal("查询组件失败")
		}
		v := comps[0].Comps[2].(*Vel)
		if v.VX != 3 || v.VY != 4 {
			t.Errorf("组件更新失败: got %+v", v)
		}
	})

	t.Run("系统更新", func(t *testing.T) {
		world := NewWorld() // 使用默认内存预算
		world.CreateEntity(
			&Pos{X: 1, Y: 1},
			&Vel{VX: 2, VY: 3},
		)
		world.Update() // 使NewEntity生效

		world.AddSys(&MoveSys{})
		world.Update() // 系统第一次Update，MoveSys生效

		comps := world.QueryComps([]int32{1})
		if len(comps) != 1 {
			t.Fatal("查询组件失败")
		}
		p := comps[0].Comps[1].(*Pos)
		if p.X != 3 || p.Y != 4 {
			t.Errorf("系统更新失败: got pos %+v", p)
		}
	})

	t.Run("移除组件", func(t *testing.T) {
		world := NewWorld() // 使用默认内存预算
		e := world.CreateEntity(
			&Pos{X: 0, Y: 0},
			&Vel{VX: 1, VY: 1},
		)
		world.Update() // 使NewEntity生效

		world.RemoveCom(e, 2) // 移除 Velocity
		world.Update()        // 使RemoveCom生效
		if world.HasCom(e, 2) {
			t.Error("组件移除失败")
		}

		comps := world.QueryComps([]int32{1})
		if len(comps) != 1 {
			t.Error("实体丢失")
		}
	})

	t.Run("销毁实体", func(t *testing.T) {
		world := NewWorld() // 使用默认内存预算
		e := world.CreateEntity(&Pos{X: 0, Y: 0})
		world.Update() // 使NewEntity生效
		world.RemoveEntity(e)
		world.Update() // 使DestroyEntity生效

		if world.HasCom(e, 1) {
			t.Error("实体销毁失败")
		}

		comps := world.QueryComps([]int32{1})
		if len(comps) != 0 {
			t.Error("销毁的实体仍然存在")
		}
	})
}

func TestECS_FullFlow(t *testing.T) {
	world := NewWorld() // 使用默认内存预算
	// 创建实体
	e1 := world.CreateEntity(&Pos{X: 10, Y: 20}, &Vel{VX: 1, VY: 2})
	e2 := world.CreateEntity(&Pos{X: 5, Y: 5})
	e3 := world.CreateEntity(&Vel{VX: 3, VY: 4})
	world.Update() // 使NewEntity生效

	// 检查组件
	if !world.HasCom(e1, 1) || !world.HasCom(e1, 2) {
		t.Error("e1组件缺失")
	}
	if !world.HasCom(e2, 1) || world.HasCom(e2, 2) {
		t.Error("e2组件状态错误")
	}
	if world.HasCom(e3, 1) || !world.HasCom(e3, 2) {
		t.Error("e3组件状态错误")
	}

	// 添加组件
	world.AddCom(e2, &Vel{VX: 2, VY: 2})
	world.Update() // 使AddCom生效
	if !world.HasCom(e2, 2) {
		t.Error("e2添加Vel失败")
	}

	// 注册系统
	world.AddSys(&MoveSys{})
	world.Update() // 系统第一次Update，MoveSys生效

	// 检查系统更新
	comps := world.QueryComps([]int32{1, 2})
	if len(comps) != 2 {
		t.Errorf("拥有Pos和Vel的实体数量错误: %d", len(comps))
	}
	for _, c := range comps {
		pos := c.Comps[1].(*Pos)
		vel := c.Comps[2].(*Vel)
		if pos.X != 11 && pos.X != 7 {
			t.Errorf("系统更新X错误: %+v", pos)
		}
		if pos.Y != 22 && pos.Y != 7 {
			t.Errorf("系统更新Y错误: %+v", pos)
		}
		if vel.VX < 1 || vel.VY < 2 {
			t.Errorf("Vel值异常: %+v", vel)
		}
	}

	// 移除组件
	world.RemoveCom(e1, 2)
	world.Update() // 使RemoveCom生效
	if world.HasCom(e1, 2) {
		t.Error("e1移除Vel失败")
	}
	if !world.HasCom(e1, 1) {
		t.Error("e1移除Vel后Pos丢失")
	}

	// 销毁实体
	world.RemoveEntity(e2)
	world.Update() // 使DestroyEntity生效
	if world.HasCom(e2, 1) || world.HasCom(e2, 2) {
		t.Error("e2销毁后组件仍存在")
	}

	// 查询剩余实体
	entities := world.Query([]int32{2})
	if len(entities) != 1 {
		t.Errorf("剩余Vel实体数量错误: %d", len(entities))
	}
}

func BenchmarkECS_Performance(b *testing.B) {
	world := NewWorld()
	numEntities := 100000
	for i := range numEntities {
		pos := &Pos{X: i, Y: i}
		vel := &Vel{VX: 1, VY: 1}
		world.CreateEntity(pos, vel)
	}
	world.AddSys(&MoveSys{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

// 完整性能测试：不同实体数量、不同组件组合、不同系统数量
type Health struct{ HP int }

func (h *Health) GetComType() int32 { return 3 }

type Tag struct{ T int }

func (t *Tag) GetComType() int32 { return 4 }

func BenchmarkECS_FullPerf(b *testing.B) {

	entityCounts := []int{1000, 10000, 30000, 50000}
	for _, n := range entityCounts {
		b.Run("Entities_"+string(rune(n)), func(b *testing.B) {
			world := NewWorld(WithNeewWorld_isSysParallel(false))
			for i := range n {
				switch i % 4 {
				case 0:
					world.CreateEntity(&Pos{X: i, Y: i}, &Vel{VX: 1, VY: 1})
				case 1:
					world.CreateEntity(&Pos{X: i, Y: i}, &Health{HP: 100})
				case 2:
					world.CreateEntity(&Vel{VX: 1, VY: 1}, &Tag{T: 1})
				case 3:
					world.CreateEntity(&Health{HP: 100}, &Tag{T: 2})
				}
			}
			world.Update()
			world.AddSys(&MoveSys{})
			world.AddSys(&HealthSys{})
			world.AddSys(&TagSys{})
			b.ResetTimer()
			for b.Loop() {
				world.Update()
			}
		})
	}
}

// 查询性能测试
func BenchmarkECS_Query(b *testing.B) {
	world := NewWorld()
	numEntities := 50000
	for i := range numEntities {
		world.CreateEntity(&Pos{X: i, Y: i}, &Vel{VX: 1, VY: 1})
	}
	world.Update()

	for b.Loop() {
		_ = world.QueryComps([]int32{1})
	}
}

// 组件增删性能测试
func BenchmarkECS_AddRemoveCom(b *testing.B) {
	world := NewWorld()
	numEntities := 10000
	entities := make([]int64, 0, numEntities)
	for i := range numEntities {
		e := world.CreateEntity(&Pos{X: i, Y: i})
		entities = append(entities, e)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			world.AddCom(e, &Vel{VX: 1, VY: 1})
			world.RemoveCom(e, 2)
		}
	}
}

// 综合测试
func TestECS_EdgeCases(t *testing.T) {
	t.Run("销毁不存在的实体", func(t *testing.T) {
		world := NewWorld()
		world.RemoveEntity(9999) // 不应 panic
		world.Update()
	})

	t.Run("HasCom不存在实体", func(t *testing.T) {
		world := NewWorld()
		world.Update()
		if world.HasCom(12345, 1) {
			t.Error("不存在实体不应有组件")
		}
	})
}
func TestECS_CreateManyEntities(t *testing.T) {
	world := NewWorld()
	entities := make([]int64, 0, 100)
	for i := range 100 {
		e := world.CreateEntity(&Pos{X: i, Y: i}, &Vel{VX: i, VY: i})
		entities = append(entities, e)
	}
	world.Update() // 使所有实体生效

	for i, e := range entities {
		if !world.HasCom(e, 1) || !world.HasCom(e, 2) {
			t.Errorf("实体%d组件缺失", i)
		}
	}
}

type CounterSys struct {
	Sys
	counter int
}

func (cs *CounterSys) Update(world *World, dt int64) {
	cs.counter++
}

func (cs *CounterSys) IsParallel() bool { return false }

func TestECS_AutoUpdate(t *testing.T) {
	world := NewWorld(WithNewWorld_AutoUpdateCD(time.Millisecond * 10)) // 自动更新间隔10ms
	sys := &CounterSys{}
	world.AddSys(sys)
	world.StartAutoUpdate()

	time.Sleep(50 * time.Millisecond)
	world.StopAutoUpdate()

	if sys.counter < 4 || sys.counter > 7 {
		t.Errorf("自动更新次数异常: got %d", sys.counter)
	}
}

// 事件测试用例
type TestEvent struct {
	Value int
}

func (e *TestEvent) GetEventType() int32 { return 1001 }
func TestECS_Event(t *testing.T) {
	world := NewWorld()
	received := 0
	world.OnEvent(1001, func(ev IEvent) {
		te := ev.(*TestEvent)
		if te.Value == 42 {
			received++
		}
	})
	world.SendEvent(&TestEvent{Value: 42})
	world.Update()
	if received != 1 {
		t.Errorf("事件未正确接收: got %d", received)
	}
}
