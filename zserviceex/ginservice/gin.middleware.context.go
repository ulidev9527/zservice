package ginservice

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 扩展 Context 中间件
func GinMiddlewareContext(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zctx := zservice.NewContext()

		zctx.RequestIP = ctx.ClientIP()
		zctx.AuthToken = ctx.Request.Header.Get(zservice.S_C2S_Token)
		zctx.ClientSign = ctx.Request.Header.Get(zservice.S_C2S_Sign)
		zctx.ClientTime = zservice.StringToUint32(ctx.Request.Header.Get(zservice.S_C2S_Time))

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
				zctx.LogErrorf("GIN %v %v %v %v %v :Q %v :E %v :ST %v",
					zctx.RequestIP, ctx.Request.Method, ctx.Request.URL,
					ctx.Writer.Status(), zctx.Since(), reqParams, e, string(buf[:stackSize]),
				)
				ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fatal, "error": zctx.TraceID})
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
