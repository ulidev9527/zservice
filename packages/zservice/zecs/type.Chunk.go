package zecs

import (
	"sync"
)

// Chunk 块结构，存储同 Archetype 的实体和组件
type Chunk struct {
	archetypeID    []int32
	entities       []int64
	comps          []map[int32]ICom
	capacity       int
	count          int
	lock           sync.RWMutex
	entityIndexMap map[int64]int // 实体ID -> 数组索引
	compsCache     []ICom
}

// 创建新 Chunk
func NewChunk(archetypeID []int32, capacity int) *Chunk {
	return &Chunk{
		archetypeID:    archetypeID,
		entities:       make([]int64, 0, capacity),
		comps:          make([]map[int32]ICom, 0, capacity),
		capacity:       capacity,
		count:          0,
		entityIndexMap: make(map[int64]int, capacity),
		compsCache:     make([]ICom, 0, capacity*len(archetypeID)),
	}
}

// 判断是否有空位
func (c *Chunk) HasSpace() bool {
	return c.count < c.capacity
}

// 添加实体及其组件
func (c *Chunk) AddEntity(entity int64, comps map[int32]ICom) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.count >= c.capacity {
		return false
	}
	idx := c.count
	if idx < len(c.entities) {
		c.entities[idx] = entity
		c.comps[idx] = comps
	} else {
		c.entities = append(c.entities, entity)
		c.comps = append(c.comps, comps)
	}
	c.entityIndexMap[entity] = idx
	// 缓存组件引用
	for _, com := range comps {
		c.compsCache = append(c.compsCache, com)
	}
	c.count++
	return true
}

// 移除实体（swap-and-pop）
func (c *Chunk) RemoveEntity(entity int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	idx, ok := c.entityIndexMap[entity]
	if !ok {
		return
	}
	lastIdx := c.count - 1
	if idx != lastIdx {
		c.entities[idx] = c.entities[lastIdx]
		c.comps[idx] = c.comps[lastIdx]
		c.entityIndexMap[c.entities[idx]] = idx
	}
	c.entities = c.entities[:lastIdx]
	c.comps = c.comps[:lastIdx]
	delete(c.entityIndexMap, entity)
	c.count--
}

// 获取实体的组件集合（无锁快速访问）
func (c *Chunk) GetEntityComps(entity int64) map[int32]ICom {
	idx, ok := c.entityIndexMap[entity]
	if ok && idx < len(c.comps) {
		return c.comps[idx]
	}
	return nil
}

// 获取实体在 chunk 中的索引
func (c *Chunk) IndexOf(entity int64) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for i, e := range c.entities {
		if int64(e) == entity {
			return i
		}
	}
	return -1
}

// 查询实体是否存在
func (c *Chunk) HasEntity(entity int64) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, exists := c.entityIndexMap[entity]
	return exists
}
