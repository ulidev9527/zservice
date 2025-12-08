package main

import (
	"os"
	"path/filepath"
	"regexp"
)

func Run_Gen_GO_regexp(dir string) error {
	re_grpc := regexp.MustCompile(`_grpc\.pb\.go$`)
	re_newXXX := regexp.MustCompile(`(in|out) :=\s*new\((.+)\)`)
	re_add___isPool := regexp.MustCompile(`type (\w+) struct \{`)
	re_msg := regexp.MustCompile(`\.pb\.go$`)

	re_reset := regexp.MustCompile(`func \(x \*.+\) Reset\(\) \{`)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !(filepath.Ext(path) == ".go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newContent := ""

		// grpc 文件 添加对象池回收
		if re_grpc.MatchString(path) {
			newContent = re_newXXX.ReplaceAllStringFunc(string(content), func(s string) string {
				m := re_newXXX.FindStringSubmatch(s)
				if len(m) == 3 {
					switch m[1] {
					case "in":
						return "in := Get_" + m[2] + "()\n	defer in.Put()"
					case "out":
						return "out := Get_" + m[2] + "()"
					}
				}
				return s
			})
		} else if re_msg.MatchString(path) {
			// pb msg 文件添加对象池验证
			newContent = re_add___isPool.ReplaceAllStringFunc(string(content), func(s string) string {
				return s + "\n\t__isInPool bool"
			})
			newContent = re_reset.ReplaceAllStringFunc(newContent, func(s string) string {
				return s + "\n\tx.reset_field()\n\treturn"
			})
		} else {
			return nil
		}

		if newContent != string(content) {
			return os.WriteFile(path, []byte(newContent), info.Mode())
		}
		return nil
	})
}
