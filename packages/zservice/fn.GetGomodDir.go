package zservice

import (
	"os"
	"path/filepath"
)

// 获取 go.mod 目录
func GetGomodDir(inputPath string) (string, *Error) {

	// 找到是否是 go.mod 目录, 不是则向上查询
	for {
		if _, err := os.Stat(filepath.Join(inputPath, "go.mod")); err != nil {
			if os.IsNotExist(err) {
				inputPath = filepath.Dir(inputPath)
			} else {
				return "", NewError("can`t find go.mod file")
			}
		} else {
			break
		}
	}
	return inputPath, nil
}
