package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_AddAsset(ctx *zservice.Context, in *zauth_pb.AddAsset_REQ) *zauth_pb.AssetInfo_RES {

	md5Str := zservice.Md5Bytes(in.FileBytes)

	savePath := fmt.Sprintf(FI_UploadDir, md5Str)
	saveDir := filepath.Dir(savePath)

	// 检查文件是否存在
	if _, e := os.Stat(saveDir); e != nil {
		if os.IsNotExist(e) {
			if e := os.MkdirAll(saveDir, 0750); e != nil {
				ctx.LogError(e)
				return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_Fail}
			}
		} else {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_Reject}
		}
	}

	// 存储文件
	if _, e := os.Stat(savePath); e != nil {
		if os.IsNotExist(e) {

			if e := os.WriteFile(savePath, in.FileBytes, 0750); e != nil {
				ctx.LogError(e)
				return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_Fail}
			}

		} else {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_Fail}
		}
	}

	// 创建资源信息
	tab, e := AssetCreate(ctx, &zauth_pb.AssetInfo{
		Name:    in.Name,
		Md5:     md5Str,
		Expires: in.Expires,
		Size:    uint64(len(in.FileBytes)),
	})
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
	} else {
		return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_SUCC, Info: &zauth_pb.AssetInfo{
			Name:    tab.Name,
			Md5:     tab.MD5,
			Token:   tab.Token,
			Expires: tab.Expires.UnixMilli(),
			Size:    tab.Size,
		}}
	}

}
