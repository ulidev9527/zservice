package zauth_ex

import (
	"net/http"
	"strings"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

// gin 检查权限中间件
// @isZauthSelf 是否是自己，自己会调用 internal 中的权限逻辑, 否则通过 grpc 调用
func GinCheckAuthMiddleware(zs *zservice.ZService, isZauthSelf ...bool) gin.HandlerFunc {

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
		zctx.AuthSign = zservice.MD5String(ctx.Request.UserAgent()) // 生成签名

		in := &zauth_pb.CheckAuth_REQ{
			Auth: string(zservice.JsonMustMarshal([]string{zservice.GetServiceName(), strings.ToLower(ctx.Request.Method), ctx.Request.URL.Path})),
		}

		res, e := func() (*zauth_pb.CheckAuth_RES, error) {
			if isSelf {
				return internal.Logic_CheckAuth(zctx, in), nil
			}
			return zauth.GetGrpcClient().CheckAuth(zctx, in)
		}()

		if e != nil {
			zctx.LogError(e)
			ctx.JSON(http.StatusOK, &zglobal.Default_RES{
				Code: zglobal.Code_Fail,
				Msg:  zctx.TraceID,
			})
			ctx.Abort()
			return
		}
		if res.Code != zglobal.Code_SUCC {

			ctx.JSON(http.StatusOK, &zglobal.Default_RES{
				Code: res.GetCode(),
				Msg:  zctx.TraceID,
			})
			ctx.Abort()
			return
		}
		if res.IsTokenRefresh {
			zctx.AuthToken = res.Token
			ginservice.SyncHeader(ctx)
		}
		zctx.UID = res.Uid

		ctx.Next()
	}
}
