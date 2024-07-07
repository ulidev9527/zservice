package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func UploadConfigAsset(ctx *zservice.Context, in *zauth_pb.UploadConfigAsset_REQ) *zauth_pb.UploadConfigAsset_RES {

	if in.Service == "" {
		in.Service = ctx.TraceService
	}

	if res, e := grpcClient.UploadConfigAsset(ctx, in); e != nil {
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
