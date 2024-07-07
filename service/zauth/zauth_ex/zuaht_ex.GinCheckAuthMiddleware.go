package zauth_ex

import (
	"net/http"
	"strings"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

// gin 检查权限中间件
// @isZauthSelf 是否是自己，自己会调用 internal 中的权限逻辑, 否则通过 grpc 调用
func GinCheckAuthMiddleware(isZauthSelf ...bool) gin.HandlerFunc {

	isSelf := false
	if len(isZauthSelf) > 0 {
		isSelf = isZauthSelf[0]
	}

	return func(ctx *gin.Context) {

		if ctx.Request.URL.Path == "" || ctx.Request.URL.Path == "/" { // 根目录不进行权限验证
			ctx.Next()
			return
		}

		zctx := ginservice.GetCtxEX(ctx)
		zctx.AuthTokenSign = zservice.MD5String(ctx.Request.UserAgent()) // 生成签名

		in := &zauth_pb.CheckAuth_REQ{
			Service:   zservice.GetServiceName(),
			Action:    strings.ToLower(ctx.Request.Method),
			Path:      ctx.Request.URL.Path,
			Token:     zctx.AuthToken,
			TokenSign: zctx.AuthTokenSign,
		}

		res, e := func() (*zauth_pb.CheckAuth_RES, error) {
			if isSelf {
				return internal.Logic_CheckAuth(zctx, in), nil
			}
			return zauth.CheckAuth(zctx, in), nil
		}()

		if e != nil {
			zctx.LogError(e)
			ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail, "msg": zctx.TraceID})
			ctx.Abort()
			return
		}

		if zctx.AuthToken != res.Token { // 刷新 token
			zctx.AuthToken = res.Token
			ginservice.SyncHeader(ctx)
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
