package internal

import (
	"strings"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zglobal"
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"
)

// 获取文件配置
func GetFileConfig(ctx *zservice.Context, in *zconfig_pb.GetFileConfig_REQ) (uint32, string) {

	if e := ParserFile(in.FileName, zglobal.E_ZConfig_Parser_Excel); e != nil {
		ctx.LogError(e)
		return e.GetCode(), ""
	}

	fKey := redisservice.FormatKey(RK_FileConfig, in.FileName)

	// 获取全部
	if in.Keys == "" {
		val, e := Redis.HGetAll(ctx, fKey).Result()
		if e != nil {
			ctx.LogError(e)
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}
		if len(val) == 0 {
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}

		return zglobal.Code_SUCC, string(zservice.JsonMustMarshal(val))
	}

	// 获取指定
	keyArr := strings.Split(in.Keys, ",")
	newArr := zservice.ListFilterString(keyArr, func(item string) bool {
		return item != ""
	})
	if len(keyArr) == 1 {
		val, e := Redis.HGet(ctx, fKey, newArr[0]).Result()
		if e != nil {
			ctx.LogError(e)
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}
		if val == "" {
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}
		return zglobal.Code_SUCC, val
	} else {
		valueArr, e := Redis.HMGet(ctx, fKey, keyArr...).Result()
		if e != nil {
			ctx.LogError(e)
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}
		if len(valueArr) == 0 {
			return zglobal.Code_Zconfig_GetConfigFail, ""
		}
		return zglobal.Code_SUCC, string(zservice.JsonMustMarshal(valueArr))
	}
}
