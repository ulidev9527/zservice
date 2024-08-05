package zservice

import (
	"fmt"
	"os"
)

// go 协程
func Go(f func()) {
	go f()
}

// 写入文件到临时目录
func WriteFileToTempDir(name string, data []byte) *Error {
	if e := os.WriteFile(fmt.Sprintf("%s/%s", os.TempDir(), name), data, 0644); e != nil {
		return NewError(e).SetCode(Code_Fatal)
	}
	return nil
}
