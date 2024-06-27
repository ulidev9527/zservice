package ginservice

import (
	"bytes"
	"fmt"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

// gin 服务扩展
type ginResWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ginResWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ginResWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 获取扩展的上下文
func GetCtxEX(ctx *gin.Context) *zservice.Context {
	z, has := ctx.Get(GIN_contextEX_Middleware_Key)
	if !has {
		return nil
	}
	zctx := z.(*zservice.Context)
	zctx.GinCtx = ctx
	return zctx
}

// 同步 header 信息
func SyncHeader(ctx *gin.Context) {
	zctx := GetCtxEX(ctx)
	if zctx.ClientSign == "" {
		zctx.ClientSign = zservice.RandomMD5()
	}
	ctx.Header(zservice.S_C2S, fmt.Sprintf("%v.%v", zctx.ContextS2S.AuthToken, zctx.ContextS2S.ClientSign))
}
