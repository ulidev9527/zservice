package internal

import (
	"fmt"
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func gin_post_config_service_uploadFileConfig(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	file, e := ctx.FormFile("file")
	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_Fail})
		return
	}
	serviceName := ctx.Param("service")
	filePath := fmt.Sprintf(FI_ServiceConfigFile, serviceName, file.Filename)
	ctx.SaveUploadedFile(file, filePath)

	ctx.JSON(http.StatusOK, Logic_ConfigSyncServiceFileConfig(zctx, &zauth_pb.ConfigSyncServiceFileConfig_REQ{
		Service:  serviceName,
		FilePath: filePath,
		Parser:   zservice.StringToUint32(ctx.PostForm("parser")),
	}))
}
