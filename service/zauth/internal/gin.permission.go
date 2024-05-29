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

	ctx.JSON(200, Logic_PermissionCreate(zctx, &zauth_pb.PermissionInfo{
		Name:    req.Name,
		Service: req.Service,
		Action:  req.Action,
		Path:    req.Path,
		State:   req.State,
	}))

}

// 获取权限
func gin_handle_get_permission(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)
	p := zservice.StringToInt32(ctx.Query("p"))
	si := zservice.StringToInt32(ctx.Query("si"))
	se := ctx.Query("se")

	ctx.JSON(200, Logic_PermissionListGet(zctx, &zauth_pb.PermissionListGet_REQ{
		Page:   p,
		Size:   si,
		Search: se,
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
