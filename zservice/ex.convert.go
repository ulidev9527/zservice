package zservice

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// 结构体转json
func Convert_StructToJson(o any) string {
	v, err := json.Marshal(o)
	if err != nil {
		println(err.Error())
		return ""
	}

	return string(v)
}

// json 转结构体
func TryConvert_JsonToStruct(jsonStr string, o interface{}) bool {
	err := json.Unmarshal([]byte(jsonStr), &o)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// json 转 结构体
func Convert_JsonToStruct(jsonStr string, o interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &o)
}

func Convert_Int64ToByte(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// int64 转字符串
func Convert_Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int 转字符串
func Convert_IntToString(i int) string {
	return strconv.Itoa(i)
}

// uint to string
func Convert_UIntToString(u uint) string {
	return fmt.Sprintf("%d", u)
}

// bool 转换字符串
func Convert_BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// string to boolean
func Convert_StringToBoolean(s string) bool {
	s = strings.ToLower(s)
	if s == "" || s == "false" || s == "0" {
		return false
	} else {
		return true
	}
}

// string to int
func Convert_StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	return i
}

// string to float32
func Convert_StringToFloat32(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		i = 0
	}
	return float32(i)
}

// string to uint
func Convert_StringToUInt(str string) uint {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}

// string to int64
func Convert_StringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

// urne to string
func Convert_RuneToString(r rune) string {
	return string(r)
}
