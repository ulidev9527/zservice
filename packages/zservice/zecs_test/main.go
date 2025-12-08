package zecs_test

import (
	"zserviceapps/packages/zservice"
	"zserviceapps/packages/zservice/zecs"
)

const (
	ECS_Com_Attack = 1
	ECS_Com_HP     = 2
)

type Com_Attk struct {
	zecs.ICom
	Atk int
}

func (c *Com_Attk) GetComType() int32 { return ECS_Com_Attack }

type Com_HP struct {
	zecs.ICom
	HP int
	cd int64
}

func (c *Com_HP) GetComType() int32 { return ECS_Com_HP }

type Sys_AutoAddHP struct {
	zecs.Sys
}

func (s *Sys_AutoAddHP) Update(world *zecs.World, dt int64) {
	world.QueryCom([]int32{ECS_Com_HP}, func(entity int64, comps map[int32]zecs.ICom) {

		com := comps[ECS_Com_HP].(*Com_HP)
		com.cd -= dt
		if com.cd > 0 {
			return
		}
		com.cd = 3000

		com.HP += 10
	})
}

func main() {

	w := zecs.NewWorld()

	for range 1_000_000 {
		w.CreateEntity(&Com_HP{})
	}

	w.StopAutoUpdate()

	zservice.WaitStop()

}
