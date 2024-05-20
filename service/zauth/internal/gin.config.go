package internal

import (
	"fmt"
	"net/http"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initConfig() {

	Gin.GET("/config", func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)

		auth := ctx.Query("auth")
		zctx.LogWarn(auth)

		ctx.String(200, "ok")
	})

	Gin.GET("/fileconfig_reset", func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)

		fileName := ctx.Query("fileName")
		if e := ParserFile(fileName, zglobal.E_ZConfig_Parser_Excel); e != nil {
			zctx.LogError(e)
			ctx.String(http.StatusOK, fmt.Sprint(e.GetCode()))
		} else {
			ctx.String(200, "ok")
		}

	})
}
