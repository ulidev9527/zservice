package zservice

import (
	"encoding/json"
)

func StructTo_MapAny(v any) map[string]any {
	return JsonMustUnmarshal_MapAny(JsonMustMarshal(v))
}

// 结构体深拷贝, 将 b 的内容拷贝给 a
func StructDeepcopy(a any, b any) *Error {
	if bt, e := json.Marshal(b); e != nil {
		return NewError(e).SetCode(Code_Fatal)
	} else if e := json.Unmarshal(bt, a); e != nil {
		return NewError(e).SetCode(Code_Fatal)
	} else {
		return nil
	}
}
