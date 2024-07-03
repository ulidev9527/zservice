package internal

import (
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_ConfigSyncServiceEnvConfig(ctx *zservice.Context, in *zauth_pb.ConfigSyncServiceEnvConfig_REQ) *zauth_pb.ConfigSyncServiceEnvConfig_RES {
	if in.Service == "" || in.FilePath == "" {
		ctx.LogError("param error")
		return &zauth_pb.ConfigSyncServiceEnvConfig_RES{Code: zglobal.Code_ParamsErr}
	}

	// 验证文件
	if e := parserFileVerify(in.FilePath); e != nil {
		ctx.LogError(e)
		return &zauth_pb.ConfigSyncServiceEnvConfig_RES{Code: e.GetCode()}
	}

	// 验证 md5
	md5Str := ""
	if str, e := zservice.Md5File(in.FilePath); e != nil {
		ctx.LogError(e)
		return &zauth_pb.ConfigSyncServiceEnvConfig_RES{Code: zglobal.Code_Zauth_config_GetFileMd5Fail}
	} else {
		md5Str = str
	}
	si := zservice.MD5String(fmt.Sprint(md5Str, time.Now().String()))
	authKey := fmt.Sprint(zservice.RandomStringLow(96), si)
	rk_auth := fmt.Sprintf(RK_Config_ServiceEnvAuth, in.Service)
	if e := Redis.Set(rk_auth, authKey).Err(); e != nil {
		if !DBService.IsNotFoundErr(e) {
			ctx.LogError(e)
			return &zauth_pb.ConfigSyncServiceEnvConfig_RES{Code: zglobal.Code_Fail}
		}
	}

	return &zauth_pb.ConfigSyncServiceEnvConfig_RES{Code: zglobal.Code_SUCC, Auth: authKey}
}
