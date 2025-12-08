package zecs

// 系统接口，所有系统需实现 Update 和 GetOrder
type ISys interface {
	OnAdd(world *World)            // 当系统被添加到世界中时执行
	Update(world *World, dt int64) // 系统更新
	GetOrder() int32               // 系统执行优先级,越小越优先
	IsParallel() bool              // 是否支持并发执行
}

type Sys struct{}

func (s *Sys) OnAdd(world *World)            {}
func (s *Sys) Update(world *World, dt int64) {}
func (s *Sys) GetOrder() int32               { return 0 }
func (s *Sys) IsParallel() bool              { return true }
