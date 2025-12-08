package zservice

import (
	"strconv"
)

// bool 转换字符串
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// bool 转 0/1 字符串
func BoolToString01(b bool) string {
	if b {
		return "1"
	} else {
		return "0"
	}
}

// bool 转 0/1 字符串
func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
func BoolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}

// urne to string
func Convert_RuneToString(r rune) string {
	return string(r)
}
