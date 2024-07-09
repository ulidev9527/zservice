package zauth

import (
	"net/http"
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

// gin 检查权限中间件
// fn 具体的检查方法
func GinCheckAuthMiddleware(s *ginservice.GinService, fn func(*zservice.Context, *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if ctx.Request.URL.Path == "" || ctx.Request.URL.Path == "/" { // 根目录不进行权限验证
			ctx.Next()
			return
		}

		zctx := s.GetCtx(ctx)
		zctx.AuthTokenSign = zservice.MD5String(ctx.Request.UserAgent()) // 生成签名

		res := fn(zctx, &zauth_pb.CheckAuth_REQ{
			Service:   zservice.GetServiceName(),
			Action:    strings.ToLower(ctx.Request.Method),
			Path:      ctx.Request.URL.Path,
			Token:     zctx.AuthToken,
			TokenSign: zctx.AuthTokenSign,
		})

		if zctx.AuthToken != res.Token { // 刷新 token
			zctx.AuthToken = res.Token
			s.SyncHeader(ctx)
		}
		zctx.UID = res.Uid

		if res.Code != zservice.Code_SUCC {
			ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail, "msg": zctx.TraceID})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
