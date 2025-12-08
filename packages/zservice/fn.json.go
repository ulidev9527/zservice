package zservice

import "encoding/json"

func JsonMustMarshal(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		LogError(e)
		return nil
	}
	return b
}
func JsonMustMarshalString(v any) string {
	return string(JsonMustMarshal(v))
}

func JsonMustUnmarshal_StringArray(v []byte) []string {
	var r []string
	e := json.Unmarshal(v, &r)
	if e != nil {
		LogError(e)
		return nil
	}
	return r
}

func JsonMustUnmarshal(v []byte) any {
	var r []string
	e := json.Unmarshal(v, &r)
	if e != nil {
		LogError(e)
		return nil
	}
	return r
}

func JsonMustUnmarshal_MapAny(v []byte) map[string]any {
	var r map[string]any
	e := json.Unmarshal(v, &r)
	if e != nil {
		LogError(e)
		return nil
	}
	return r
}
