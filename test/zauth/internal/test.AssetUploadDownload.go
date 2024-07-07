package internal

import (
	"fmt"
	"os"
	"time"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Test_AssetUploadDownload() *zservice.ZService {
	ctx := zservice.NewContext()
	return zservice.NewService("asset up test", func(z *zservice.ZService) {
		defer z.StartDone()

		filePath := "Config_Item.xlsx"
		bt, e := os.ReadFile(fmt.Sprint("static/", filePath))
		if e != nil {
			ctx.LogError(e)
			return
		}

		upRES := zauth.UploadAsset(ctx, &zauth_pb.UploadAsset_REQ{
			Name:    "Config_Item.xlsx",
			Expires: time.Now().Add(zservice.Time_1m).UnixMilli(),
			Data:    bt,
		})
		if upRES.Code != zservice.Code_SUCC {
			ctx.LogError("UploadAsset Fail:", upRES)
			return
		}

		downRES := zauth.DownloadAsset(ctx, &zauth_pb.DownloadAsset_REQ{
			AssetID: upRES.Info.AssetID,
		})
		if downRES.Code != zservice.Code_SUCC {
			ctx.LogError("DownloadAsset Fail:", upRES)
			return
		}

		ctx.LogInfo("DownloadAsset Success:", downRES)

		if e := os.WriteFile(fmt.Sprintf("static/down.%s", downRES.Info.Name), downRES.Info.Data, 0750); e != nil {
			ctx.LogError(e)
			return
		}
	})
}
