package zecs

import "testing"

func TestWorld_GetCom(t *testing.T) {
	world := NewWorld()
	pos := &Pos{X: 10, Y: 20}
	vel := &Vel{VX: 1, VY: 2}
	e := world.CreateEntity(pos, vel)
	world.Update()

	// 正常获取
	gotPos := world.GetCom(e, 1)
	if gotPos == nil {
		t.Fatal("GetCom 未获取到 Pos")
	}
	if p, ok := gotPos.(*Pos); !ok || p.X != 10 || p.Y != 20 {
		t.Errorf("GetCom 获取到的 Pos 错误: %+v", gotPos)
	}

	gotVel := world.GetCom(e, 2)
	if gotVel == nil {
		t.Fatal("GetCom 未获取到 Vel")
	}
	if v, ok := gotVel.(*Vel); !ok || v.VX != 1 || v.VY != 2 {
		t.Errorf("GetCom 获取到的 Vel 错误: %+v", gotVel)
	}

	// 获取不存在的组件
	gotNil := world.GetCom(e, 999)
	if gotNil != nil {
		t.Error("GetCom 获取不存在组件应为nil")
	}

	// 获取不存在实体的组件
	gotNil2 := world.GetCom(99999, 1)
	if gotNil2 != nil {
		t.Error("GetCom 获取不存在实体应为nil")
	}
}
