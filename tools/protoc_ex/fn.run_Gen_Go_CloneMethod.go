package main

import (
	"fmt"
	"strings"
)

// 生成 go Clone 相关方法
func run_Gen_Go_CloneMethod(outputStr string, protoMessageInfos []*ProtoMessageInfo) string {

	fileStrMap := map[string]string{}
	for _, info := range protoMessageInfos {
		str := ""
		if _, ok := fileStrMap[info.FileName]; ok {
			str = fileStrMap[info.FileName]
		}

		fieldsStr := "\n"
		for _, field := range info.Fields {

			switch field.MType {
			case "[]*":
				fieldsStr += fmt.Sprintf("    for _, v := range x.%s { if v != nil { clone.AddOne%s(v.Clone()) } }\n", field.Name, field.Name)
			case "*":
				fieldsStr += fmt.Sprintf("    if x.%s != nil { clone.Set%s(x.%s.Clone()) }\n", field.Name, field.Name, field.Name)
			default:
				fieldsStr += fmt.Sprintf("    clone.Set%s(x.%s)\n", field.Name, field.Name)
			}

		}

		cloneStr := `
// 克隆 ${type_name}
func (x *${type_name}) Clone() *${type_name} {
    clone := Get_${type_name}()
${fields}
    return clone
}
`

		cloneStr = strings.ReplaceAll(cloneStr, "${type_name}", info.Name)
		cloneStr = strings.ReplaceAll(cloneStr, "${fields}", fieldsStr)

		str += cloneStr
		// 存储
		fileStrMap[info.FileName] = str
	}

	// 将代码输出到文件
	for fileName, str := range fileStrMap {
		outputStr += fmt.Sprint("\n\n// ---------------------------------------- run_Gen_Go_CloneMethod => 👇 FileName: ", fileName, "\n")
		outputStr += str
		outputStr += fmt.Sprint("// ---------------------------------------- run_Gen_Go_CloneMethod => 👆 FileName: ", fileName, "\n\n")

	}
	return outputStr

}
