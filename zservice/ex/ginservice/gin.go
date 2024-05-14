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

func (grw *ginResWriter) Write(b []byte) (int, error) {
	grw.body.Write(b)
	return grw.ResponseWriter.Write(b)
}

// 获取扩展的上下文
func GetCtxEX(ctx *gin.Context) *zservice.Context {
	z, has := ctx.Get(GIN_contextEX_Middleware_Key)
	if !has {
		return nil
	}
	return z.(*zservice.Context)
}

type GinService struct {
	*zservice.ZService
	Ginengine *gin.Engine
}

type GinServiceConfig struct {
	Name string // 服务名
	Addr string // 监听地址

	OnStart func(*gin.Engine) // 启动的回调
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// gin 服务扩展
func NewGinService(c *GinServiceConfig) *GinService {

	if c == nil {
		zservice.LogPanic("GinServiceConfig is nil")
		return nil
	}
	name := "GinService"
	if c.Name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	gs := &GinService{}
	gs.Ginengine = gin.New()

	// 服务
	gs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		if c.OnStart != nil {
			c.OnStart(gs.Ginengine)
		}

		go func() {
			gs.LogInfof("ginService listen on %v", c.Addr)
			e := gs.Ginengine.Run(c.Addr)
			if e != nil {
				s.LogPanic(e)
			}
		}()
		go func() {
			s.StartDone()
		}()

	})

	// 中间件
	gs.Ginengine.Use(GinCORSMiddleware(gs.ZService))
	gs.Ginengine.Use(GinContextEXMiddleware(gs.ZService))
	gs.Ginengine.Use(GinAuthEXMiddleware(gs.ZService))

	return gs
}
