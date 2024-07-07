package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func UploadAsset(ctx *zservice.Context, in *zauth_pb.UploadAsset_REQ) *zauth_pb.AssetInfo_RES {
	if res, e := grpcClient.UploadAsset(ctx, in); e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_500Err}
	} else {
		return res
	}
}
