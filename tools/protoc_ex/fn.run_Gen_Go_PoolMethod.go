package main

import (
	"fmt"
	"strings"
)

// 生成 go 对象池 代码
func run_Gen_Go_PoolMethod(outputStr string, protoMessageInfos []*ProtoMessageInfo) string {

	fileStrMap := map[string]string{}

	for _, info := range protoMessageInfos {

		str := ""
		if _, ok := fileStrMap[info.FileName]; ok {
			str = fileStrMap[info.FileName]
		}

		// 字段重置字符串生成
		allResetValue := "\n"

		for _, field := range info.Fields {

			reset_value := ""
			ex_reset_value := ""

			switch field.MType {
			case "string":
				reset_value = "\"\""
			case "int", "int32", "int64", "uint32", "uint64", "float32":
				reset_value = "0"
			case "bool":
				reset_value = "false"
			case "[]*":
				reset_value = "nil"
				putName := field.Type[3:]
				ex_reset_value = fmt.Sprintf(" for _, item := range msg.%s { put_%s(item) }", field.Name, putName)

			case "[]":
				reset_value = "nil"
			case "*":
				putName := field.Type[1:]
				reset_value = "nil"
				ex_reset_value = fmt.Sprintf("put_%s(msg.%s)", putName, field.Name)
			case "E_", "EConst":
				reset_value = "0"
			default:
				reset_value = fmt.Sprintf("// %s__%s__???", field.Name, field.Type)
			}

			allResetValue += fmt.Sprintf("    %s\n", ex_reset_value)
			if strings.HasPrefix(reset_value, "//") {
				allResetValue += fmt.Sprintf("    // 无法处理:msg.%s=%s\n", field.Name, reset_value)
			} else {
				allResetValue += fmt.Sprintf("    msg.%s=%s\n", field.Name, reset_value)
			}

		}

		tp := `
// ${type_name} 消息池
var pool_${type_name} = &sync.Pool{
    New: func() any {
        msg:= &${type_name}{}
        return msg
    },
}

// 回收 ${type_name} 消息
func (x *${type_name}) Put() {
	// put_${type_name}(x)
}

// 回收 ${type_name} 消息 并返回数据字节内容
func (x *${type_name}) Put_Bytes() []byte {
	// defer put_${type_name}(x)
	return x.ToBytes()
}

// 回收 ${type_name} 消息 并返回数据 json 内容
func (x *${type_name}) Put_Json() string {
	// defer put_${type_name}(x)
	return x.ToJson()
}

// 从对象池中获取 ${type_name} 消息
func Get_${type_name}() *${type_name} {
    msg := pool_${type_name}.Get().(*${type_name})
	// msg.__isInPool = false
	return msg
}

// 回收 ${type_name} 消息
func put_${type_name}(msg *${type_name}) {
    // if msg == nil {
    //     return
    // }

	// if msg.__isInPool {
    //     zservice.LogErrorCallerf(3, "%s is In Pool", reflect.TypeOf(msg).String())
	// 	return
	// }

    // msg.Reset()
    
	// msg.__isInPool = true

    // pool_${type_name}.Put(msg)
}

// 重置 ${type_name} 内容
func (msg *${type_name}) reset_field() {
    ${allResetValue}
}
	
// 转换为 []byte
func (x *${type_name}) ToBytes() []byte {
	return zservice.ProtobufMustMarshal(x)
}

// 转换为 json 字符串
func (x *${type_name}) ToJson() string {
	return zservice.JsonMustMarshalString(x)
}`

		tp = strings.ReplaceAll(tp, "${allResetValue}", allResetValue)
		tp = strings.ReplaceAll(tp, "${type_name}", info.Name)

		str += tp

		// 存储
		fileStrMap[info.FileName] = str

	}

	// 将代码输出到文件
	for fileName, str := range fileStrMap {
		outputStr += fmt.Sprint("\n\n// ---------------------------------------- run_Gen_Go_PoolMethod => 👇 FileName: ", fileName, "\n")
		outputStr += str
		outputStr += fmt.Sprint("// ---------------------------------------- run_Gen_Go_PoolMethod => 👆 FileName: ", fileName, "\n\n")

	}
	return outputStr
}
