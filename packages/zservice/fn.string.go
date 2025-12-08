package zservice

import (
	"fmt"
	"strconv"
	"strings"
)

// 格式化
func StringSprint(v ...any) string {
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
	return strconv.Quote(StringSprint(v...))
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

func StringSplitUint32(s string, sep string, clearEmpty ...bool) []uint32 {
	arr := StringSplit(s, sep, clearEmpty...)
	newArr := []uint32{}
	for _, v := range arr {
		newArr = append(newArr, StringToUint32(v))
	}

	return newArr
}
func StringSplitInt32(s string, sep string, clearEmpty ...bool) []int32 {
	arr := StringSplit(s, sep, clearEmpty...)
	newArr := []int32{}
	for _, v := range arr {
		newArr = append(newArr, StringToInt32(v))
	}

	return newArr
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
func StringToUint8(s string) uint8 {
	i, err := strconv.ParseUint(strings.TrimSpace(s), 10, 8)
	if err != nil {
		return 0
	}
	return uint8(i)
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
	i, err := strconv.ParseInt(strings.TrimSpace(str), 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

// string to float32
func StringToFloat32(str string) float32 {
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 32)
	if err != nil {
		i = 0
	}
	return float32(i)
}

// string to uint
func StringToUint(str string) uint {
	i, err := strconv.ParseUint(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}

// string to uint32
func StringToUint32(str string) uint32 {
	i, err := strconv.ParseUint(strings.TrimSpace(str), 10, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}

// string to int64
func StringToInt64(str string) int64 {
	i, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

// string to uint64
func StringToUint64(str string) uint64 {
	i, err := strconv.ParseUint(strings.TrimSpace(str), 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

// string to float32
func StringToFloat(str string) float32 {
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 32)
	if err != nil {
		i = 0
	}
	return float32(i)
}

// string to float64
func StringToFloat64(str string) float64 {
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		i = 0
	}
	return i
}

// string to []float64
func StringToSplitFloat64(str string) float64 {
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		i = 0
	}
	return i
}

func StringSplitFloat64(s string, sep string, clearEmpty ...bool) []float64 {
	arr := StringSplit(s, sep, clearEmpty...)
	newArr := []float64{}
	for _, v := range arr {
		newArr = append(newArr, StringToFloat64(v))
	}

	return newArr
}

// 秒转毫秒时间戳
func StringTo_Time_S_TO_MS(s string) int64 {
	return int64(StringToFloat64(s) * 1000)
}
