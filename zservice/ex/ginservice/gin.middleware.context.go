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

		zctx := func() *zservice.Context { // 提取 S2S / C2S 信息
			s2sStr := ctx.Request.Header.Get(zservice.S_S2S)
			c2sStr := ctx.Request.Header.Get(zservice.S_C2S)
			zzctx := zservice.NewContext(s2sStr)
			zzctx.ContextS2S.RequestIP = ctx.ClientIP()
			if c2sStr == "" {
				return zzctx
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

		grw := &ginResWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = grw

		reqParams := ""

		switch strings.Split(ctx.Request.Header.Get("Content-Type"), ";")[0] {
		case "application/json": // 处理 json 类型数据
			reqBody, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// gogin数据读取一次后无法再次读取，所以需要重新写入一份
			dst := &bytes.Buffer{}
			if len(reqBody) > 0 {
				if e := json.Compact(dst, reqBody); e != nil {
					zctx.LogError(e)
				} else {
					reqParams = dst.String()
				}

			}
		}

		defer func() {
			//放在匿名函数里,e捕获到错误信息，并且输出
			e := recover()
			if e != nil {
				buf := make([]byte, 1<<12)
				stackSize := runtime.Stack(buf, true)
				zctx.LogErrorf("GIN %v %v %v %v %v :Q %v :E %v %v",
					zctx.RequestIP, ctx.Request.Method, ctx.Request.URL,
					ctx.Writer.Status(), zctx.Since(), reqParams, e, string(buf[:stackSize]),
				)
				ctx.JSON(200, gin.H{"code": 0, "error": zctx.TraceID})
			}
		}()

		ctx.Next()

		// 打印日志
		zctx.LogInfof("GIN %v %v %v %v %v :Q %v :S %v",
			zctx.RequestIP, ctx.Request.Method, ctx.Request.URL,
			ctx.Writer.Status(), zctx.Since(),
			reqParams, grw.body.String(),
		)
	}
}
