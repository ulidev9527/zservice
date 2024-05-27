package ginservice

import (
	"bytes"
	"encoding/json"
	"io"
	"runtime"
	"strings"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 扩展 Context 中间件
func GinMiddlewareContext(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		zctx := func() *zservice.Context { // 提取处 C2S 信息
			zzctx := zservice.NewEmptyContext()
			c2sStr := ctx.Request.Header.Get(zservice.S_C2S)
			if c2sStr == "" {
				return zservice.NewEmptyContext()
			}

			if len(c2sStr) == 65 {
				c2sArr := strings.Split(c2sStr, ".")
				if len(c2sArr) == 2 {
					zzctx.ContextS2S.AuthToken = c2sArr[0]
					zzctx.ContextS2S.ClientSign = c2sArr[1]
				}
			}

			return zzctx
		}()

		ctx.Set(GIN_contextEX_Middleware_Key, zctx)

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
				buf := make([]byte, 1<<10)
				stackSize := runtime.Stack(buf, true)
				zctx.LogErrorf("GIN %v %v %v %v %v :Q %v :E %v %v",
					ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL,
					ctx.Writer.Status(), zctx.Since(), reqParams, e, string(buf[:stackSize]),
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