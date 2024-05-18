package zservice

func StructTo_MapAny(v any) map[string]any {
	return JsonMustUnmarshal_MapAny(JsonMustMarshal(v))
}
