package zservice

import "slices"

// 数组去重
func ListRemoveDuplicates(arr interface{}) interface{} {
	a := map[interface{}]int{}
	r := []interface{}{}

	for _, v := range arr.([]interface{}) {
		if a[v] == 0 {
			a[v] = 1
			r = append(r, v)
		}
	}
	return r
}

// 移除数组中重头开始的指定元素
func ListRemoveFirst(arr interface{}, i int) interface{} {
	return append(arr.([]interface{})[:0], arr.([]interface{})[i:]...)
}

// 过滤器
func ListFilter(arr any, fn func(item any) bool) any {
	newArr := &[]any{}
	for _, item := range arr.([]any) {
		if fn(item) {
			*newArr = append(*newArr, item)
		}
	}

	return *newArr
}

// 过滤器
func ListFilterString(arr []string, fn func(item string) bool) []string {
	newArr := []string{}
	for _, item := range arr {
		if fn(item) {
			newArr = append(newArr, item)
		}
	}
	return newArr
}

// 数组中是否有某个值
func ListHas(arr interface{}, fn func(item any) bool) bool {
	for _, v := range arr.([]interface{}) {
		if fn(v) {
			return true
		}
	}
	return false
}

// ui32 arr 去重
func List_UI32_Duplicates(list []uint32) []uint32 {
	m := map[uint32]int{}
	for i := 0; i < len(list); i++ {
		if _, ok := m[list[i]]; ok {
			list = slices.Delete(list, i, i+1)
			i--
		} else {
			m[list[i]] = 0
		}
	}
	m = nil
	return list
}

// 复制数组
func List_Clone[T any](src []T) []T {
	if src == nil {
		return nil // 处理 nil 输入
	}
	dst := make([]T, len(src)) // 预分配目标切片（容量=长度，避免额外内存）
	copy(dst, src)             // 使用内置 copy 函数（高效，底层优化）
	return dst
}
