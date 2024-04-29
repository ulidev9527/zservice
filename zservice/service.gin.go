package zservice

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// gin 服务扩展
type ginResWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (grw *ginResWriter) Write(b []byte) (int, error) {
	grw.body.Write(b)
	return grw.ResponseWriter.Write(b)
}

// 中间件
// CORS跨域中间件
func GinCORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}

var __gin_contextEX_Middleware_Key = "__gin_contextEX_Middleware_Key"

// 扩展 Context 中间件
func GinContextEXTMiddleware(zs *ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zctx := NewContext(zs, ctx.Request.Header.Get(S_TraceKey))
		ctx.Set(__gin_contextEX_Middleware_Key, ctx)

		reqParams := ""
		bodyStr := ""

		switch strings.Split(ctx.Request.Header.Get("Content-Type"), ";")[0] {
		case "application/json": // 处理 json 类型数据
			reqBody, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// gogin数据读取一次后无法再次读取，所以需要重新写入一份
			dst := &bytes.Buffer{}
			if e := json.Compact(dst, reqBody); e != nil {
				zctx.LogError(e)
			} else {
				reqParams = dst.String()
			}

			grw := &ginResWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: ctx.Writer,
			}
			ctx.Writer = grw
			bodyStr = grw.body.String()
		}

		ctx.Next()

		zctx.LogInfof("GIN %v %v %v %v %v REQ %v RES %v",
			ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL,
			ctx.Writer.Status(), zctx.Since(),
			reqParams, bodyStr,
		)
	}
}

// 获取扩展的上下文
func GetGinCtxEX(ctx *gin.Context) *ZContext {
	z, has := ctx.Get(__gin_contextEX_Middleware_Key)
	if !has {
		return nil
	}
	return z.(*ZContext)
}
