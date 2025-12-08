package main

import (
	"fmt"
	"strings"
)

// 生成 go Set 相关方法
func run_Gen_Go_FieldMethod(outputStr string, protoMessageInfos []*ProtoMessageInfo) string {

	fileStrMap := map[string]string{}

	for _, info := range protoMessageInfos {
		str := ""
		if _, ok := fileStrMap[info.FileName]; ok {
			str = fileStrMap[info.FileName]
		}
		for _, field := range info.Fields {

			tpStr := `
// 设置 ${field_name}.${field_type} 字段值
func (x *${type_name}) Set${field_name}(v ${field_type})*${type_name} {
    x.${field_name} = v
    return x
}
`

			tpStr = strings.ReplaceAll(tpStr, "${type_name}", info.Name)
			tpStr = strings.ReplaceAll(tpStr, "${field_name}", field.Name)
			tpStr = strings.ReplaceAll(tpStr, "${field_type}", field.Type)

			str += tpStr

		}

		// 存储
		fileStrMap[info.FileName] = str

	}

	// 将代码输出到文件
	for fileName, str := range fileStrMap {
		outputStr += fmt.Sprint("\n\n// ---------------------------------------- run_Gen_Go_FieldMethod => 👇 FileName: ", fileName, "\n")
		outputStr += str
		outputStr += fmt.Sprint("// ---------------------------------------- run_Gen_Go_FieldMethod => 👆 FileName: ", fileName, "\n\n")

	}
	return outputStr
}
