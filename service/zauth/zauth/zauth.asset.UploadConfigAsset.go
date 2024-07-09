package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func UploadConfigAsset(ctx *zservice.Context, in *zauth_pb.UploadConfigAsset_REQ) *zauth_pb.UploadConfigAsset_RES {

	if in.Service == "" {
		ctx.LogError("param in nil")
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_ParamsErr}
	}

	if res, e := grpcClient.UploadConfigAsset(ctx, in); e != nil {
		ctx.LogError(e)
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
