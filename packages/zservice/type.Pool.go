package zservice

import (
	"sync"
)

// Pool 是一个泛型对象池
type Pool[T any] struct {
	pool     sync.Pool
	createFn func() T
	putFN    func(T)
}

// NewPool 创建一个新的对象池
// create: 创建新对象的函数
// put: 回收时调用
// maxSize: 池中最大对象数量，<=0表示无限制
func NewPool[T any](createFN func() T, putFN func(T)) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return createFN()
			},
		},
		createFn: createFN,
		putFN:    putFN,
	}
}

// Get 从池中获取一个对象
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put 将对象放回池中
func (p *Pool[T]) Put(obj T) {
	p.putFN(obj)
	p.pool.Put(obj)
}
