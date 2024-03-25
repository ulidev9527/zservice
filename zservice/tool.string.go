package zservice

import "fmt"

// 格式化
func Sprint(v ...any) string {
	l := len(v)
	if l <= 1 {
		return fmt.Sprint(v...)
	}
	nv := make([]any, 0)
	for _, vv := range v {
		nv = append(nv, vv, " ")
	}
	return fmt.Sprint(nv[:len(nv)-1]...)
}
