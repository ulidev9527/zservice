package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initGinConfig() {

	Gin.POST("/config/:service/uploadFileConfig", func(ctx *gin.Context) {

		zctx := ginservice.GetCtxEX(ctx)
		file, e := ctx.FormFile("file")
		if e != nil {
			zctx.LogError(e)
			ctx.JSON(200, gin.H{"code": zglobal.Code_ErrorBreakoff})
			return
		}
		serviceName := ctx.Param("service")
		filePath := fmt.Sprintf(FI_ServiceConfigFile, serviceName, file.Filename)
		ctx.SaveUploadedFile(file, filePath)

		ctx.JSON(200, Logic_ConfigSyncServiceFileConfig(zctx, &zauth_pb.ConfigSyncServiceFileConfig_REQ{
			Service:  serviceName,
			FilePath: filePath,
			Parser:   zservice.StringToUint32(ctx.PostForm("parser")),
		}))

	})

	Gin.POST("/config/:service/uploadEnvConfig", func(ctx *gin.Context) {

		zctx := ginservice.GetCtxEX(ctx)
		file, e := ctx.FormFile("file")
		if e != nil {
			zctx.LogError(e)
			ctx.JSON(200, gin.H{"code": zglobal.Code_ErrorBreakoff})
			return
		}
		serviceName := ctx.Param("service")
		filePath := fmt.Sprintf(FI_ServiceEnvFile, serviceName)
		ctx.SaveUploadedFile(file, filePath)

		ctx.JSON(200, Logic_ConfigSyncServiceEnvConfig(zctx, &zauth_pb.ConfigSyncServiceEnvConfig_REQ{
			Service:  serviceName,
			FilePath: filePath,
		}))
	})

	Gin.GET("/config/serviceEnvConfig/:auth", func(ctx *gin.Context) {

		zctx := ginservice.GetCtxEX(ctx)
		auth := ctx.Param("auth")

		ctx.JSON(200, Logic_ConfigGetServiceEnvConfig(zctx, &zauth_pb.ConfigGetServiceEnvConfig_REQ{
			Service: zctx.TraceService,
			Auth:    auth,
		}))

	})
}
