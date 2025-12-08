package zecs

import (
	"sync"
	"unsafe"
)

const (
	MinChunkSize = 8    // 最小chunk大小
	MaxChunkSize = 1024 // 最大chunk大小
)

// Archetype 原型结构
type Archetype struct {
	id           []int32  // 组件类型集合
	chunks       []*Chunk // 所有块
	memoryBudget int      // 替换chunkSize为memoryBudget
	chunkSize    int      // 实际计算出的chunk大小
	lock         sync.RWMutex
}

// 创建新 Archetype
func NewArchetype(id []int32, memoryBudget int) *Archetype {
	return &Archetype{
		id:           id,
		chunks:       []*Chunk{},
		memoryBudget: memoryBudget,
		chunkSize:    0, // 将在首次AddEntity时计算
	}
}

func (a *Archetype) calculateChunkSize(comps map[int32]ICom) int {
	// 基础内存开销
	entitySize := int(unsafe.Sizeof(int64(0)))
	mapEntrySize := int(unsafe.Sizeof(map[int32]ICom{}))

	// 计算组件实际内存占用
	totalCompSize := 0
	for _, comp := range comps {
		// 使用unsafe.Sizeof计算实际组件大小
		totalCompSize += int(unsafe.Sizeof(comp))
	}

	// 计算每个实体的总内存开销 (包括对齐填充)
	totalEntitySize := entitySize + mapEntrySize + totalCompSize
	alignedSize := (totalEntitySize + 7) &^ 7 // 8字节对齐

	// 根据内存预算计算chunk大小
	chunkSize := a.memoryBudget / alignedSize

	// 确保在合理范围内
	if chunkSize < MinChunkSize {
		chunkSize = MinChunkSize
	} else if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	return chunkSize
}

// 分配实体到 Chunk
func (a *Archetype) AddEntity(entity int64, comps map[int32]ICom) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	// 首次添加实体时计算chunkSize
	if a.chunkSize == 0 {
		a.chunkSize = a.calculateChunkSize(comps)
	}

	// 查找有空位的 Chunk
	for _, chunk := range a.chunks {
		if chunk.HasSpace() {
			return chunk.AddEntity(entity, comps)
		}
	}

	// 没有空位，新建 Chunk
	newChunk := NewChunk(a.id, a.chunkSize)
	a.chunks = append(a.chunks, newChunk)
	return newChunk.AddEntity(entity, comps)
}

// 查询实体所在 Chunk
func (a *Archetype) FindChunk(entity int64) *Chunk {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, chunk := range a.chunks {
		if chunk.HasEntity(entity) {
			return chunk
		}
	}
	return nil
}

// 移除实体
func (a *Archetype) RemoveEntity(entity int64) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, chunk := range a.chunks {
		if chunk.HasEntity(entity) {
			chunk.RemoveEntity(entity)
			break
		}
	}
}

// 并发遍历 Archetype 下所有实体及其组件
func (a *Archetype) Each(fn func(entity int64, comps map[int32]ICom)) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, chunk := range a.chunks {
		chunk.lock.RLock()
		for i := 0; i < chunk.count; i++ {
			fn(chunk.entities[i], chunk.comps[i])
		}
		chunk.lock.RUnlock()
	}
}
