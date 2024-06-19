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

// int to uint32
func IntToUint32(i int) uint32 {
	return uint32(i)
}

// uint to string
func UintToString(u uint) string {
	return fmt.Sprintf("%d", u)
}

// uint32 to string
func Uint32ToString(u uint32) string {
	return fmt.Sprintf("%d", u)
}

// uint32 to int
func Uint32ToInt(u uint32) int {
	return int(u)
}

// uint64 to string
func UInt64ToString(u uint64) string {
	return fmt.Sprintf("%d", u)
}
