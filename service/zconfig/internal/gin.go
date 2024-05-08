package internal

import (
	"zservice/internal/httpservice"

	"github.com/gin-gonic/gin"
)

var Gin *gin.Engine

func InitGin() {
	Gin.GET("/config", func(ctx *gin.Context) {
		zctx := httpservice.GetGinCtxEX(ctx)

		auth := ctx.Query("auth")
		zctx.LogWarn(auth)

		ctx.String(200, "ok")
	})
}
