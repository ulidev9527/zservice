package zservice

import (
	"fmt"
	"strconv"
	"strings"
)

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

// 去掉换行的格式化
func SprintQuote(v ...any) string {
	return strconv.Quote(Sprint(v...))
}

func StringSplit(s string, sep string, clearEmpty ...bool) []string {
	arr := strings.Split(s, sep)
	if len(clearEmpty) > 0 && clearEmpty[0] {
		arr = ListFilterString(arr, func(item string) bool {
			return item != ""
		})
	}
	return arr
}
