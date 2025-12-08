package extgo

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strings"
	"zserviceapps/packages/zservice"

	"github.com/spf13/viper"
)

type ParamInfo struct {
	Name string
	Type string
}

func Run(conf *viper.Viper, inputFile string) {
	var inputBodyStr string
	if f, e := os.ReadFile(inputFile); e != nil {
		zservice.LogError(inputFile, e)
		return
	} else {
		inputBodyStr = string(f)
	}

	structMap := parseStruct(inputFile)
	if structMap == nil {
		return
	}

	appendString := ""

	for _, structInfo := range structMap {

		has_method_Pack := false

		// 有 pack 方法
		if method_Pack := structInfo.GetMethod("Pack"); method_Pack.Name != "" {
			if len(method_Pack.Params) == 1 {
				if method_Pack.Params[0].Type == "*flatbuffers.Builder" {
					has_method_Pack = true
				}
			}
		}

		// 有 Table 方法
		has_method_Table := false
		if m := structInfo.GetMethod("Table"); m.Name != "" {
			if len(m.Returns) == 1 && m.Returns[0].Type == "flatbuffers.Table" {
				has_method_Table = true
			}
		}

		// json 结构体给原始结构体添加方法
		sourceStructName := strings.TrimSuffix(structInfo.Name, "T") // 原始结构体
		if sourceStructInfo, has := structMap[sourceStructName]; has {
			if methodUnPack := sourceStructInfo.GetMethod("UnPackTo"); methodUnPack.Name != "" {
				for _, fieldInfo := range structInfo.Fields {

					if strings.HasPrefix(fieldInfo.Type, "[]") {
						if strings.HasPrefix(fieldInfo.Type, "[]*") {

						} else if strings.HasPrefix(fieldInfo.Type, "[]byte") {

						} else {
							appendString += addFieldValueTypeArr(sourceStructInfo, fieldInfo)
						}
					}

				}
			}
		}

		if has_method_Pack {
			appendString += addFBTMethod(structInfo)
		}

		if has_method_Table {
			appendString += addFBMethod(structInfo)
		}

	}

	if appendString != "" {
		inputBodyStr += "\n"
		inputBodyStr += appendString

		// 替换import
		inputBodyStr = strings.ReplaceAll(inputBodyStr, `flatbuffers "github.com/google/flatbuffers/go"`, `
	"zserviceapps/packages/zservice"

	flatbuffers "github.com/google/flatbuffers/go"
		`)

		// 添加字段 tag
		regex := regexp.MustCompile("(`json:(\".+\"))")
		inputBodyStr = regex.ReplaceAllString(inputBodyStr, "$1 mapstructure:$2")

		// 添加崩溃监控
		regex = regexp.MustCompile(`(func \(rcv \*.+\{)`)
		inputBodyStr = regex.ReplaceAllString(inputBodyStr, "$1\n    defer zservice.RecoverFix()")
	}

	// 重新输出
	if e := os.WriteFile(inputFile, []byte(inputBodyStr), 0644); e != nil {
		zservice.LogError(inputFile, e)
	}

}

func addFBTMethod(structInfo StructInfo) string {

	str := `
func (t *${structName}) Reset() {
${resetStr}
}

func (t *${structName}) ToJson()string{	return  zservice.JsonMustMarshalString(t) }

// 转换为 FB 对象
func (t *${structName}) To${structNameFB}() *${structNameFB} {
	bts := t.ToBytes()
	return GetRootAs${structNameFB}(bts, 0)
}

// 转换为 FB 对象并且立即回收自己
func (t *${structName}) To${structNameFB}_Put() *${structNameFB} {
	bts := t.ToBytes()
	PutPool_${structName}(t)
	return GetRootAs${structNameFB}(bts, 0)
}

func (t *${structName}) ToBytes() []byte {
	builder := GetBuilder()
	builder.Finish(t.Pack(builder))
	buf := builder.FinishedBytes()
	newBuf := make([]byte, len(buf))
	copy(newBuf, buf)
	PutBuilder(builder)
	buf = nil
	return newBuf
}

// 转换为 bytes 并回收
func (t *${structName}) ToBytes_Put() []byte {
	newBuf := t.ToBytes()
	PutPool_${structName}(t)
	return newBuf
}

var pool_${structName} = zservice.NewPool(func() *${structName} { return &${structName}{} },func(t *${structName}) { t.Reset() })

func GetPool_${structName}() *${structName} {
	return pool_${structName}.Get()
}
func PutPool_${structName}(t *${structName}) {
	pool_${structName}.Put(t)
}

func (t *${structName}) Clone() *${structName} {
	buf := t.ToBytes()
	cmd := GetRootAs${structNameFB}(buf, 0)
	t2 := GetPool_${structName}()
	cmd.UnPackTo(t2)
	return t2
}
	`

	// 重置字符串
	str = strings.ReplaceAll(str, "${resetStr}", func() string {
		resetStr := ""
		for _, fieldInfo := range structInfo.Fields {
			switch fieldInfo.Type {
			case "string":
				resetStr += fmt.Sprintf("    t.%s = \"\"\n", fieldInfo.Name)
			case "int", "int32", "int64", "uint32", "uint64", "float32":
				resetStr += fmt.Sprintf("    t.%s = 0\n", fieldInfo.Name)
			default:
				if strings.HasPrefix(fieldInfo.Type, "[]") {
					resetStr += fmt.Sprintf("    t.%s = nil\n", fieldInfo.Name)
				} else {
					zservice.LogError("can`t use ", fieldInfo.Name, fieldInfo.Type)
				}

			}
		}
		return resetStr
	}())

	// 添加链式方法
	for _, fieldInfo := range structInfo.Fields {
		switch fieldInfo.Type {
		case "string", "int", "int32", "int64", "uint32", "uint64", "float32":
			s := `
func (t *${structName}) Set${fieldName}(val ${fieldType}) *${structName} {
	t.${fieldName} = val
	return t
}`
			s = strings.ReplaceAll(s, "${fieldName}", fieldInfo.Name)
			s = strings.ReplaceAll(s, "${fieldType}", fieldInfo.Type)
			str += s

		default:
			if strings.HasPrefix(fieldInfo.Type, "[]") {
				s := `
func (t *${structName}) Set${fieldName}(val ${fieldType}) *${structName} {
	t.${fieldName} = val
	return t
}`
				s = strings.ReplaceAll(s, "${fieldName}", fieldInfo.Name)
				s = strings.ReplaceAll(s, "${fieldType}", fieldInfo.Type)
				str += s
			} else {
				zservice.LogError("can`t use ", fieldInfo.Name, fieldInfo.Type)
			}
		}
	}

	str = strings.ReplaceAll(str, "${structName}", structInfo.Name)
	str = strings.ReplaceAll(str, "${structNameFB}", strings.TrimSuffix(structInfo.Name, "T"))

	return str
}

func addFBMethod(structInfo StructInfo) string {
	str := `
func (rcv *${structName}) Clone() *${structName} {
	return GetRootAs${structName}(zservice.List_Clone(rcv._tab.Bytes), 0)
}
`
	str = strings.ReplaceAll(str, "${structName}", structInfo.Name)
	return str
}

// 添加数组方法
func addFieldValueTypeArr(structInfo StructInfo, fieldInfo FieldInfo) string {
	str := `
func (rcv *${structName}) ${fieldName}_ToList() ${fieldType} {
	kLen := rcv.${fieldName}Length()
	switch kLen {
	case 0:
		return nil
	case 1:
		return ${fieldType}{rcv.${fieldName}(0)}
	default:
		list := make(${fieldType}, kLen)
		for i := range kLen {
			list = append(list, rcv.${fieldName}(i))
		}
		return list
	}
}
`
	str = strings.ReplaceAll(str, "${structName}", structInfo.Name)
	str = strings.ReplaceAll(str, "${fieldName}", fieldInfo.Name)
	str = strings.ReplaceAll(str, "${fieldType}", fieldInfo.Type)
	return str
}

// 解析结构体
func parseStruct(inputFile string) map[string]StructInfo {

	structMap := make(map[string]StructInfo)

	// 解析 Go 文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return nil
	}

	// 第一次遍历：收集所有结构体定义
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						structName := typeSpec.Name.Name
						structMap[structName] = StructInfo{
							Name:    structName,
							Fields:  []FieldInfo{},
							Methods: []MethodInfo{},
						}
					}
				}
			}
		}
	}

	// 第二次遍历：收集结构体字段和方法
	for _, decl := range node.Decls {
		// 处理结构体字段
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structName := typeSpec.Name.Name
						if _, exists := structMap[structName]; !exists {
							continue
						}

						// 更新字段信息
						fields := make([]FieldInfo, 0)
						for _, field := range structType.Fields.List {
							fieldType := getTypeName(field.Type)
							for _, name := range field.Names {
								fields = append(fields, FieldInfo{
									Name: name.Name,
									Type: fieldType,
								})
							}
							// 处理匿名字段
							if len(field.Names) == 0 {
								fields = append(fields, FieldInfo{
									Name: fieldType, // 匿名字段的类型名作为字段名
									Type: fieldType,
								})
							}
						}

						// 更新结构体信息
						info := structMap[structName]
						info.Fields = fields
						structMap[structName] = info
					}
				}
			}
		}

		// 处理方法
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
			// 获取接收器类型
			recvType := getReceiverType(funcDecl.Recv)
			if recvType == "" {
				continue
			}

			// 检查是否是已知结构体的方法
			if _, exists := structMap[recvType]; !exists {
				continue
			}

			// 处理方法参数
			params := make([]ParamInfo, 0)
			if funcDecl.Type.Params != nil {
				for _, param := range funcDecl.Type.Params.List {
					paramType := getTypeName(param.Type)
					if len(param.Names) > 0 {
						for _, name := range param.Names {
							params = append(params, ParamInfo{
								Name: name.Name,
								Type: paramType,
							})
						}
					} else {
						params = append(params, ParamInfo{
							Name: "",
							Type: paramType,
						})
					}
				}
			}

			// 处理返回参数
			returns := make([]ParamInfo, 0)
			if funcDecl.Type.Results != nil {
				for _, result := range funcDecl.Type.Results.List {
					returnType := getTypeName(result.Type)
					if len(result.Names) > 0 {
						for _, name := range result.Names {
							returns = append(returns, ParamInfo{
								Name: name.Name,
								Type: returnType,
							})
						}
					} else {
						returns = append(returns, ParamInfo{
							Name: "",
							Type: returnType,
						})
					}
				}
			}

			method := MethodInfo{
				Name:    funcDecl.Name.Name,
				Params:  params,
				Returns: returns,
			}

			// 更新结构体信息
			info := structMap[recvType]
			info.Methods = append(info.Methods, method)
			structMap[recvType] = info
		}
	}

	return structMap
}

// 获取接收器类型
func getReceiverType(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}

	field := recv.List[0]
	switch t := field.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// 获取类型名称
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeName(t.X)
	case *ast.ArrayType:
		return "[]" + getTypeName(t.Elt)
	case *ast.SelectorExpr:
		return getTypeName(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + getTypeName(t.Key) + "]" + getTypeName(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return "func" + getFuncTypeSignature(t)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + getTypeName(t.Value)
		case ast.RECV:
			return "<-chan " + getTypeName(t.Value)
		default:
			return "chan " + getTypeName(t.Value)
		}
	case *ast.Ellipsis:
		return "..." + getTypeName(t.Elt)
	default:
		return reflect.TypeOf(t).String()
	}
}

// 获取函数类型签名
func getFuncTypeSignature(funcType *ast.FuncType) string {
	var params, results []string

	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			typeName := getTypeName(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					params = append(params, name.Name+" "+typeName)
				}
			} else {
				params = append(params, typeName)
			}
		}
	}

	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			typeName := getTypeName(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					results = append(results, name.Name+" "+typeName)
				}
			} else {
				results = append(results, typeName)
			}
		}
	}

	signature := "(" + strings.Join(params, ", ") + ")"
	if len(results) > 0 {
		if len(results) > 1 || results[0] != "" {
			signature += " (" + strings.Join(results, ", ") + ")"
		} else {
			signature += " " + results[0]
		}
	}

	return signature
}
