package ginservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 扩展 Context 中间件
func GinMiddlewareContext(zs *zservice.ZService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := GetContext(c)

		logStr := ""

		defer func() {
			//放在匿名函数里,e捕获到错误信息，并且输出
			e := recover()
			if e != nil {
				buf := make([]byte, 1<<12)
				stackSize := runtime.Stack(buf, true)
				ctx.LogErrorf("%v :E %v :ST %v",
					logStr, e, string(buf[:stackSize]),
				)
				c.JSON(http.StatusOK, DefResJson(zservice.Code_Fatal))
			}
		}()

		switch strings.Split(c.Request.Header.Get("Content-Type"), ";")[0] {
		case "application/json": // 处理 json 类型数据
			grw := &ginResWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}
			reqParams := ""
			c.Writer = grw
			reqBody, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// gogin数据读取一次后无法再次读取，所以需要重新写入一份
			dst := &bytes.Buffer{}
			if len(reqBody) > 0 {
				if e := json.Compact(dst, reqBody); e != nil {
					ctx.LogError(e)
				} else {
					reqParams = dst.String()
				}

			}
			logStr = fmt.Sprintf("GIN %v %v %v %v %v :Q %v",
				ctx.RequestIP, c.Request.Method, c.Request.URL,
				c.Writer.Status(), ctx.Since(),
				reqParams,
			)

		default:
			logStr = fmt.Sprintf("GIN %v %v %v %v %v",
				ctx.RequestIP, c.Request.Method, c.Request.URL,
				c.Writer.Status(), ctx.Since(),
			)
		}

		c.Next()

		switch c.Writer.Header().Get("Content-Type") {
		case "application/json": // 处理 json 类型数据
			ctx.LogInfof("%v :S %v", logStr, c.Writer.(*ginResWriter).body.String())
		default:
			ctx.LogInfo(logStr)
		}
	}
}

// 获取 gin 携带的上下文
func GetContext(c *gin.Context) *zservice.Context {
	z, has := c.Get(GIN_contextEX_Middleware_Key)
	if !has {
		z = zservice.NewContext()
		ctx := z.(*zservice.Context)
		ctx.RequestIP = c.ClientIP()
		ctx.AuthToken = c.Request.Header.Get(zservice.S_C2S_Token)
		ctx.ClientSign = c.Request.Header.Get(zservice.S_C2S_Sign)
		ctx.ClientTime = zservice.StringToUint32(c.Request.Header.Get(zservice.S_C2S_Time))

		c.Set(GIN_contextEX_Middleware_Key, ctx)
	}

	return z.(*zservice.Context)
}

// 同步请求头的 context 信息
func SyncHeader(c *gin.Context) {
	ctx := GetContext(c)

	if ctx.ClientSign != "" {
		c.Header(zservice.S_C2S_Sign, ctx.ClientSign)
		c.Header(zservice.S_C2S_Time, zservice.Int64ToString(time.Now().UnixMilli()))
	}

	c.Header(zservice.S_C2S_Token, ctx.AuthToken)
}
