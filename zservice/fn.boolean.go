package zservice

import (
	"strconv"
)

// bool 转换字符串
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// urne to string
func Convert_RuneToString(r rune) string {
	return string(r)
}
