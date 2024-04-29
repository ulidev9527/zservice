package zservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RestService struct {
	*ZService             // 服务
	gineng    *gin.Engine // gogin 对象
}

type GinResWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var (
	S_TraceHeader = "ZSERVER-TRACE"
	S_zserverCTX  = "ZSERVER_CTX"
)

func (grw *GinResWriter) Write(b []byte) (int, error) {
	grw.body.Write(b)
	return grw.ResponseWriter.Write(b)
}

// REST 配置
type ZServiceRESTConfig struct {
	Name          string             // 服务名称
	Addr          string             // 监听地址 (IP:端口 *.*.*.*:***) (:端口 :***)
	OnBeforeStart func(*RestService) // 启动前的回调
}

// 创建一个 rest 服务
func NewRestService(c *ZServiceRESTConfig) *RestService {
	if c == nil {
		LogError("ZServiceRESTConfig is nil")
		return nil
	}
	name := c.Name
	if name == "" {
		name = "DEF"
	}

	// gin 全局配置
	gin.SetMode(gin.ReleaseMode)

	// 初始化内部参数
	rest := &RestService{}
	eng := gin.New()
	zs := NewService(&ZServiceConfig{
		Name: fmt.Sprint("RestService-", name),
		OnBeforeStart: func(s *ZService) {
			if c.OnBeforeStart != nil {
				c.OnBeforeStart(rest)
			}
		},
		OnStart: func(s *ZService) {
			s.StartDone()
			e := eng.Run(c.Addr)
			if e != nil {
				s.LogError(e)
			}
		},
	})

	rest.gineng = eng
	rest.ZService = zs

	// 中间件
	eng.Use(func(ctx *gin.Context) {
		// 处理跨域
		func() {
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
		}()

		// 添加本次请求的 Context / Trace / 日志打印
		func() {

			zctx := NewContext(rest.ZService, ctx.Request.Header.Get(S_TraceHeader))
			ctx.Set(S_zserverCTX, ctx)

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

				grw := &GinResWriter{
					body:           bytes.NewBufferString(""),
					ResponseWriter: ctx.Writer,
				}
				ctx.Writer = grw
				bodyStr = grw.body.String()
			}

			ctx.Next()

			zctx.LogInfof("REST %v %v %v %v %v REQ %v RES %v",
				ctx.ClientIP(), ctx.Request.Method, ctx.Request.URL,
				ctx.Writer.Status(), zctx.Since(),
				reqParams, bodyStr,
			)
		}()
	})

	return rest
}

func (rs *RestService) AddRouter(httpMethod string, relativePath string, handlerFunc gin.HandlerFunc) {
	rs.gineng.POST("", func(ctx *gin.Context) {

	})
}
