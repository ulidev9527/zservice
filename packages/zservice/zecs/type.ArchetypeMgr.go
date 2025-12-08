package zecs

import (
	"hash/fnv"
	"slices"
	"sync"
)

type ArchetypeMgr struct {
	archetypes   map[uint64]*Archetype
	memoryBudget int // 替换chunkSize为memoryBudget
	lock         sync.RWMutex
}

func NewArchetypeMgr(memoryBudget int) *ArchetypeMgr {
	return &ArchetypeMgr{
		archetypes:   make(map[uint64]*Archetype),
		memoryBudget: memoryBudget,
	}
}

func hashComTypes(ids []int32) uint64 {
	copyIds := make([]int32, len(ids))
	copy(copyIds, ids)
	slices.Sort(copyIds)
	uniqIds := make([]int32, 0, len(copyIds))
	for i, id := range copyIds {
		if i == 0 || id != copyIds[i-1] {
			uniqIds = append(uniqIds, id)
		}
	}
	h := fnv.New64a()
	for _, id := range uniqIds {
		b := []byte{
			byte(id >> 24), byte(id >> 16), byte(id >> 8), byte(id),
		}
		h.Write(b)
	}
	return h.Sum64()
}

func (mgr *ArchetypeMgr) GetOrCreateArchetype(ids []int32) *Archetype {
	key := hashComTypes(ids)
	mgr.lock.RLock()
	a, ok := mgr.archetypes[key]
	mgr.lock.RUnlock()
	if ok {
		return a
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	if a, ok := mgr.archetypes[key]; ok {
		return a
	}
	a = NewArchetype(ids, mgr.memoryBudget)
	mgr.archetypes[key] = a
	return a
}

// 返回所有包含 typeIDs 的 Archetype
func (mgr *ArchetypeMgr) MatchArchetypes(typeIDs []int32) []*Archetype {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	// 预分配空间避免扩容
	result := make([]*Archetype, 0, len(mgr.archetypes)/2)

	// 对短typeIDs使用简单的比较
	if len(typeIDs) <= 3 {
		for _, arch := range mgr.archetypes {
			hasAll := true
			for _, tid := range typeIDs {
				if !slices.Contains(arch.id, tid) {
					hasAll = false
					break
				}
			}
			if hasAll {
				result = append(result, arch)
			}
		}
		return result
	}

	// 对长typeIDs使用map查找
	typeIDSet := make(map[int32]struct{}, len(typeIDs))
	for _, id := range typeIDs {
		typeIDSet[id] = struct{}{}
	}

	for _, arch := range mgr.archetypes {
		hasAll := true
		for id := range typeIDSet {
			if !slices.Contains(arch.id, id) {
				hasAll = false
				break
			}
		}
		if hasAll {
			result = append(result, arch)
		}
	}
	return result
}
