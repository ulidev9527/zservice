package zservice

import "encoding/json"

func JsonMustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func JsonMustUnmarshalStringArray(v []byte) []string {
	var r []string
	json.Unmarshal(v, &r)
	return r
}
