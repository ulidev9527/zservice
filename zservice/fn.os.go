package zservice

import (
	"fmt"
	"os"

	"github.com/panjf2000/ants/v2"
)

// go 原生协程
func Go(f func()) {
	go f()
}

// go ants 协程
func GO_ants(f func()) {
	e := ants.Submit(f)
	if e != nil {
		LogError(e)
	}
}

// 写入文件到临时目录
func WriteFileToTempDir(name string, data []byte) *Error {
	if e := os.WriteFile(fmt.Sprintf("%s/%s", os.TempDir(), name), data, 0644); e != nil {
		return NewError(e).SetCode(Code_Fatal)
	}
	return nil
}
