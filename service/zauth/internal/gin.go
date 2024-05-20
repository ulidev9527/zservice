package internal

import (
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

var Gin *gin.Engine

func InitGin() {
	initConfig()

	Gin.GET("/config", func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)

		auth := ctx.Query("auth")
		zctx.LogWarn(auth)

		ctx.String(200, "ok")
	})

	Gin.GET("/version", func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)
		zctx.LogPanic("eeeeeee")
	})
}
