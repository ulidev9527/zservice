package zservice

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// 参数签名编码
func _paramSignEncode(data map[string]interface{}, aid string, cv string, sv string) string {
	t := data["t"].(string)
	k := data["k"].(string)

	// 参数排序
	str := string(JsonMustMarshal(data))

	str = strings.ReplaceAll(str, "{}", "")
	str = strings.ReplaceAll(str, "[]", "")

	strArr := []rune(str)
	sort.Slice(strArr, func(i, j int) bool {
		return strArr[i] > strArr[j]
	})
	str = string(strArr)
	str = MD5String(t + k + str + aid + cv)
	str = MD5String(sv + str)
	return str
}

type T_ParamSignArgs struct {
	AID string `json:"aid"` // appid
	CV  string `json:"cv"`  // 客户端版本
	SV  string `json:"sv"`  // 签名ID
}

// 参数签名
func ParamSign(o any, psa T_ParamSignArgs) any {
	jObj := gjson.Parse(string(JsonMustMarshal(o)))
	useData := make(map[string]interface{})

	// 赋值
	jObj.ForEach(func(key, value gjson.Result) bool {
		k := key.Str
		useData[k] = value.Value()
		return true
	})

	useData["t"] = fmt.Sprint(time.Now().UnixMilli())
	delete(useData, "s")

	useData["s"] = _paramSignEncode(useData, psa.AID, psa.CV, psa.SV)
	return useData
}

// 参数签名验证
func ParamSignVerify(jsonStr string, psa T_ParamSignArgs) bool {

	jObj := gjson.Parse(jsonStr)
	useData := make(map[string]interface{})

	jObj.ForEach(func(key, value gjson.Result) bool {
		useData[key.Str] = value.Value()
		return true
	})

	s := useData["s"].(string)
	delete(useData, "s")

	return s == _paramSignEncode(useData, psa.AID, psa.CV, psa.SV)
}
