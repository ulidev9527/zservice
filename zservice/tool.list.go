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
