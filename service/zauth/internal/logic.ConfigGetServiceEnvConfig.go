package internal

import (
	"fmt"
	"io"
	"os"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"
)

func Logic_ConfigGetServiceEnvConfig(ctx *zservice.Context, in *zauth_pb.ConfigGetServiceEnvConfig_REQ) *zauth_pb.ConfigGetServiceEnvConfig_RES {
	if len(in.Auth) != 128 && ctx.TraceService == "" {
		return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_ParamsErr}
	}

	rk_auth := fmt.Sprintf(RK_Config_ServiceEnvAuth, ctx.TraceService)

	if str, e := Redis.Get(rk_auth).Result(); e != nil {
		if !redisservice.IsNilErr(e) {
			ctx.LogError(e)
			return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_ErrorBreakoff}
		}
	} else if str != in.Auth {
		return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_NotFound}
	}

	if fi, e := os.Open(fmt.Sprintf(FI_ServiceEnvFile, ctx.TraceService)); e != nil {
		ctx.LogError(e)
		return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_OpenFileErr}
	} else {
		defer fi.Close()
		if bt, e := io.ReadAll(fi); e != nil {
			return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_OpenFileErr}
		} else {
			return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_SUCC, Value: string(bt)}
		}
	}

}
