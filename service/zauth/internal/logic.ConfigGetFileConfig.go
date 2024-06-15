package internal

import (
	"fmt"
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取文件配置
func Logic_ConfigGetFileConfig(ctx *zservice.Context, in *zauth_pb.ConfigGetFileConfig_REQ) *zauth_pb.ConfigGetFileConfig_RES {
	fKey := fmt.Sprintf(RK_Config_ServiceFileConfig, in.Service, in.FileName) // 缓存 key
	has, e := Redis.Exists(fKey).Result()
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
	}
	if has == 0 {
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
	}

	// 获取全部
	if in.Keys == "" {
		val, e := Redis.HGetAll(fKey).Result()
		if e != nil {
			ctx.LogError(e)
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}
		if len(val) == 0 {
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}

		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_SUCC, Value: string(zservice.JsonMustMarshal(val))}
	}

	// 获取指定
	keyArr := strings.Split(in.Keys, ",")
	newArr := zservice.ListFilterString(keyArr, func(item string) bool {
		return item != ""
	})
	if len(keyArr) == 1 {
		val, e := Redis.HGet(fKey, newArr[0]).Result()
		if e != nil || val == "" {
			ctx.LogError(e)
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}
		if val == "" {
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_SUCC, Value: val}
	} else {
		valueArr, e := Redis.HMGet(fKey, keyArr...).Result()
		if e != nil {
			ctx.LogError(e)
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}
		if len(valueArr) == 0 {
			return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
		}
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_SUCC, Value: string(zservice.JsonMustMarshal(valueArr))}
	}
}