package ginservice

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"strings"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

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
func GinContextEXTMiddleware(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zctx := zservice.NewContext(zs, ctx.Request.Header.Get(zservice.S_TraceKey))
		ctx.Set(__gin_contextEX_Middleware_Key, zctx)

		var grw *ginResWriter
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

			grw = &ginResWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: ctx.Writer,
			}
			ctx.Writer = grw
		}

		defer func() {
			//放在匿名函数里,e捕获到错误信息，并且输出
			e := recover()
			if e != nil {
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, true)
				zctx.LogError(e, string(buf[:stackSize]))
				zctx.LogErrorf("GIN %v %v %v %v %v :Q %v ",
					ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL,
					ctx.Writer.Status(), zctx.Since(), reqParams,
				)
				ctx.String(500, "ERROR: %v", zctx.TraceID)
			}
		}()

		ctx.Next()

		if grw != nil && grw.body != nil {
			bodyStr = grw.body.String()
		}

		// 打印日志
		zctx.LogInfof("GIN %v %v %v %v %v :Q %v :S %v",
			ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL,
			ctx.Writer.Status(), zctx.Since(),
			reqParams, bodyStr,
		)
	}
}
