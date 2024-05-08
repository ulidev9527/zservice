package internal

import (
	"zservice/internal/ginservice"

	"github.com/gin-gonic/gin"
)

var Gin *gin.Engine

func InitGin() {
	Gin.GET("/config", func(ctx *gin.Context) {
		zctx := ginservice.GetGinCtxEX(ctx)

		auth := ctx.Query("auth")
		zctx.LogWarn(auth)

		ctx.String(200, "ok")
	})

	Gin.GET("/version", func(ctx *gin.Context) {
		zctx := ginservice.GetGinCtxEX(ctx)
		zctx.LogPanic("eeeeeee")
	})
}
