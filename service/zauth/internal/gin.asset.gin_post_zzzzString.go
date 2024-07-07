package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_post_zzzzString(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)

	file, e := ctx.FormFile("file")
	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	if bt, e := ginservice.ReadUploadFile(file); e != nil { // 读取文件
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	} else if res := Logic_UploadAsset(zctx, &zauth_pb.UploadAsset_REQ{ // 上传文件
		Name: file.Filename,
		Data: bt,
	}); res.Code != zservice.Code_SUCC {
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	} else if res := Logic_SetServiceKV(zctx, &zauth_pb.SetServiceKV_REQ{ // 设置服务KV
		Service: zservice.GetServiceName(),
		Key:     KV_ZZZZString,
		Value:   res.Info.AssetID,
	}); res.Code != zservice.Code_SUCC {
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
	} else if e := ZZZZString.Reload(zctx); e != nil { // 重新加载
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": e.GetCode()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_SUCC})

}
