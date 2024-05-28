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

}

// 获取权限
func gin_handle_get_permission(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)
	p := zservice.StringToInt32(ctx.Query("p"))
	si := zservice.StringToInt32(ctx.Query("si"))
	se := ctx.Query("se")

	res := Logic_GetPermissionList(zctx, &zauth_pb.GetPermissionList_REQ{
		Page:   p,
		Size:   si,
		Search: se,
	})

	if res.Code == zglobal.Code_SUCC {
		ctx.JSON(200, gin.H{"code": res.Code, "list": res.List})
	} else {
		ctx.JSON(200, gin.H{"code": res.Code})
	}

}
