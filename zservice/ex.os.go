package zservice

import "fmt"

// go 协程
func Go(f func()) {
	go f()
}

// 获取一个临时文件路径
func GetTempFilepath() string {
	return fmt.Sprint("__tmp__/", RandomMD5())
}
