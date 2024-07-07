package zservice

import "strconv"

// 是否是数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// 是否是整数
func IsInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
