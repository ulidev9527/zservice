package zservice

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

func Int64ToByte(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// int64 转字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int 转字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// uint to string
func UIntToString(u uint) string {
	return fmt.Sprintf("%d", u)
}
