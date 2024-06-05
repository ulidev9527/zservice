package internal

import (
	"net/http"
	"time"
	"zservice/service/zlog/zlog_pb"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

var GinService *ginservice.GinService
var Gin *gin.Engine

func InitGin() {
	Gin.POST("/addKV", func(ctx *gin.Context) {

		zctx := ginservice.GetCtxEX(ctx)

		req := &zlog_pb.LogKV_REQ{}

		if e := ctx.ShouldBindJSON(req); e != nil {
			zctx.LogError(e)
		} else {
			req.SaveTime = time.Now().UnixMilli()
		}

		Logic_AddLogKV(zctx, req)

		ctx.String(http.StatusOK, "")
	})
}
