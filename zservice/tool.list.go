package zservice

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
