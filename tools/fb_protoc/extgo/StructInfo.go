package extgo

// 定义存储结构
type StructInfo struct {
	Name    string
	Fields  []FieldInfo
	Methods []MethodInfo
}

func (o StructInfo) GetMethod(name string) MethodInfo {
	for _, v := range o.Methods {
		if v.Name == name {
			return v
		}
	}
	return MethodInfo{}
}
