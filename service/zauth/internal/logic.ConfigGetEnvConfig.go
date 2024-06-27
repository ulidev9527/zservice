package internal

import (
	"fmt"
	"io"
	"os"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"zservice/service/zauth/zauth_pb"
)

func Logic_ConfigGetEnvConfig(ctx *zservice.Context, in *zauth_pb.ConfigGetEnvConfig_REQ) *zauth_pb.ConfigGetServiceEnvConfig_RES {
	if in.Service == "" {
		ctx.LogError("param error")
		return &zauth_pb.ConfigGetServiceEnvConfig_RES{Code: zglobal.Code_ParamsErr}
	}

	if fi, e := os.Open(fmt.Sprintf(FI_ServiceEnvFile, in.Service)); e != nil {
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
