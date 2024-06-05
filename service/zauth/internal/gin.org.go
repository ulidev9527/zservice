package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initGinOrg() {
	Gin.GET("/org", gin_orgListGet)
	Gin.POST("/org", gin_orgCreate)
	Gin.PUT("/org", gin_orgUpdate)
}

// 创建组织
func gin_orgCreate(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)

	req := &zauth_pb.OrgInfo{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	ctx.JSON(http.StatusOK, Logic_OrgCreate(zctx, req))

}

func gin_orgListGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, Logic_OrgListGet(ginservice.GetCtxEX(ctx), &zauth_pb.OrgListGet_REQ{
		Page:   zservice.StringToUint32(ctx.Query("p")),
		Size:   zservice.StringToUint32(ctx.Query("si")),
		Search: ctx.Query("se"),
	}))
}

func gin_orgUpdate(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	req := &zauth_pb.OrgInfo{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	ctx.JSON(http.StatusOK, Logic_OrgUpdate(zctx, req))
}
