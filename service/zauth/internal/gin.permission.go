package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initGinPermission() {
	Gin.GET("/permission", gin_handle_get_permission)
	Gin.PUT("/permission", gin_handle_put_permission)
	Gin.POST("/permission", gin_handle_post_permission)
	Gin.POST("/permission/bind", gin_handle_permissionBind)

}

// 创建权限
func gin_handle_post_permission(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	req := &zauth_pb.PermissionInfo{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	ctx.JSON(200, Logic_PermissionCreate(zctx, req))

}

// 获取权限
func gin_handle_get_permission(ctx *gin.Context) {

	ctx.JSON(200, Logic_PermissionListGet(ginservice.GetCtxEX(ctx), &zauth_pb.PermissionListGet_REQ{
		Page:   zservice.StringToUint32(ctx.Query("p")),
		Size:   zservice.StringToUint32(ctx.Query("si")),
		Search: ctx.Query("se"),
	}))

}

// 修改权限
func gin_handle_put_permission(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	req := &zauth_pb.PermissionInfo{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	ctx.JSON(200, Logic_PermissionUpdate(zctx, req))
}

// 权限绑定
func gin_handle_permissionBind(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := &zauth_pb.PermissionBind_REQ{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	ctx.JSON(200, Logic_PermissionBind(zctx, req))
}
