package main

import (
	"fmt"
	"strings"
)

// 生成 go Set 相关方法
func run_Gen_Go_AddMethod(outputStr string, protoMessageInfos []*ProtoMessageInfo) string {
	fileStrMap := map[string]*strings.Builder{}

	// 模板定义
	const (
		tmplPtr = `
// 添加一个 ${field_type}
func (x *${type_name}) AddOne${field_name}(v *${field_type}) *${type_name} {
    x.${field_name} = append(x.${field_name}, v)
    return x
}
`
		tmplVal = `
// 添加一个 ${field_type}
func (x *${type_name}) AddOne${field_name}(v ${field_type}) *${type_name} {
    x.${field_name} = append(x.${field_name}, v)
    return x
}
`
	)

	for _, info := range protoMessageInfos {
		builder, ok := fileStrMap[info.FileName]
		if !ok {
			builder = &strings.Builder{}
			fileStrMap[info.FileName] = builder
		}

		for _, field := range info.Fields {
			var tmpl string
			switch field.MType {
			case "[]*":
				tmpl = tmplPtr
			case "[]":
				tmpl = tmplVal
			default:
				continue
			}
			method := strings.ReplaceAll(tmpl, "${type_name}", info.Name)
			method = strings.ReplaceAll(method, "${field_name}", field.Name)
			method = strings.ReplaceAll(method, "${field_type}", field.KType)
			builder.WriteString(method)
		}
	}

	// 输出到 outputStr
	for fileName, builder := range fileStrMap {
		outputStr += fmt.Sprintf("\n\n// ---------------------------------------- run_Gen_Go_AddMethod => 👇 FileName: %s\n", fileName)
		outputStr += builder.String()
		outputStr += fmt.Sprintf("// ---------------------------------------- run_Gen_Go_AddMethod => 👆 FileName: %s\n\n", fileName)
	}
	return outputStr
}
