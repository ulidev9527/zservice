package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_post_configAsset(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	file, e := ctx.FormFile("file")
	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	if file.Size > 10485760 { // 10MB
		zctx.LogError("file size too big:", file.Size)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Reject})
		return
	}

	service := ctx.PostForm("service")
	if service == "" {
		zctx.LogError("service is empty")
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_ParamsErr})
		return
	}

	if bt, e := ginservice.ReadUploadFile(file); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	} else {
		ctx.JSON(http.StatusOK, Logic_UploadConfigAsset(zctx, &zauth_pb.UploadConfigAsset_REQ{
			Service: service,
			Name:    file.Filename,
			Parser:  zservice.StringToUint32(ctx.PostForm("parser")),
			Data:    bt,
		}))
	}
}
