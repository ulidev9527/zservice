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

// string to boolean
func StringToBoolean(s string) bool {
	s = strings.ToLower(s)
	if s == "" || s == "false" || s == "0" {
		return false
	} else {
		return true
	}
}

// string to int, err return 0
func StringToInt(str string) int {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return i
}

// string to int32
func StringToInt32(str string) int32 {
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

// string to float32
func StringToFloat32(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		i = 0
	}
	return float32(i)
}

// string to uint
func StringToUint(str string) uint {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}

// string to int64
func StringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}
